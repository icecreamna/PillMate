package handlers

import (
	"time"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

const ymdLayout = "2006-01-02"

// ===================================================================
//                   Query/List: ดึงรายการ NotiItem
// ===================================================================

type ListNotiItemsFilter struct {
	PatientID    *uint
	MyMedicineID *uint
	GroupID      *uint
	NotiInfoID   *uint
	DateFrom     *string // "YYYY-MM-DD"
	DateTo       *string // "YYYY-MM-DD"
	TakenStatus  *bool
	NotifyStatus *bool
}

type NotiItemWithNames struct {
	ID        uint `json:"id"`
	PatientID uint `json:"patient_id"`

	MyMedicineID uint   `json:"my_medicine_id"`
	NotifyDate   string `json:"notify_date"`
	NotifyTime   string `json:"notify_time"`
	NotiInfoID   uint   `json:"noti_info_id"`

	GroupID       *uint  `json:"group_id,omitempty"`
	GroupName     string `json:"group_name,omitempty"`
	MedName       string `json:"med_name"`
	AmountPerTime string `json:"amount_per_time"`
	FormID        uint   `json:"form_id"`
	UnitID        *uint  `json:"unit_id,omitempty"`
	InstructionID *uint  `json:"instruction_id,omitempty"`

	TakenStatus  bool       `json:"taken_status"`
	TakenTimeAt  *time.Time `json:"taken_time_at,omitempty"`
	NotifyStatus bool       `json:"notify_status"`
	HasSymptom   bool       `json:"has_symptom"`

	FormName        string  `json:"form_name"`
	UnitName        *string `json:"unit_name,omitempty"`
	InstructionName *string `json:"instruction_name,omitempty"`
}

func ListNotiItems(db *gorm.DB, patientID uint, f ListNotiItemsFilter) ([]NotiItemWithNames, error) {
	q := db.Table("noti_items ni").
		Joins("LEFT JOIN forms f ON ni.form_id = f.id").
		Joins("LEFT JOIN units u ON ni.unit_id = u.id").
		Joins("LEFT JOIN instructions i ON ni.instruction_id = i.id").
		Select(`
    	ni.id,
    	ni.patient_id,
    	ni.my_medicine_id,
    	ni.group_id,
    	ni.group_name,
    	ni.noti_info_id,
    	to_char(ni.notify_date, 'YYYY-MM-DD') AS notify_date,
    	to_char(ni.notify_time AT TIME ZONE 'UTC', 'HH24:MI') AS notify_time, -- ✅ บังคับใช้ UTC
    	ni.taken_status,
    	ni.notify_status,
    	ni.med_name,
    	ni.amount_per_time,
    	ni.form_id,
    	f.form_name,
    	ni.unit_id,
    	u.unit_name,
    	ni.instruction_id,
    	i.instruction_name
		`).
		Where("ni.patient_id = ?", patientID).Where("ni.deleted_at IS NULL")

	if f.MyMedicineID != nil {
		q = q.Where("my_medicine_id = ?", *f.MyMedicineID)
	}
	if f.GroupID != nil {
		q = q.Where("group_id = ?", *f.GroupID)
	}
	if f.NotiInfoID != nil {
		q = q.Where("noti_info_id = ?", *f.NotiInfoID)
	}
	if f.TakenStatus != nil {
		q = q.Where("taken_status = ?", *f.TakenStatus)
	}
	if f.NotifyStatus != nil {
		q = q.Where("notify_status = ?", *f.NotifyStatus)
	}
	if f.DateFrom != nil && *f.DateFrom != "" {
		if parsedFrom, err := time.ParseInLocation(ymdLayout, *f.DateFrom, time.Local); err == nil {
			q = q.Where("notify_date >= ?", parsedFrom)
		}
	}
	if f.DateTo != nil && *f.DateTo != "" {
		if parsedTo, err := time.ParseInLocation(ymdLayout, *f.DateTo, time.Local); err == nil {
			q = q.Where("notify_date <= ?", parsedTo)
		}
	}

	var rows []struct {
		ID              uint
		PatientID       uint
		MyMedicineID    uint
		GroupID         *uint
		GroupName       string
		NotiInfoID      uint
		NotifyDate      string
		NotifyTime      string
		TakenStatus     bool
		NotifyStatus    bool
		MedName         string
		AmountPerTime   string
		FormID          uint
		FormName        string
		UnitID          *uint
		UnitName        *string
		InstructionID   *uint
		InstructionName *string
	}

	if err := q.Order("ni.notify_date, ni.notify_time").Scan(&rows).Error; err != nil {
		return nil, err
	}

	// ✅ map ไปยัง struct สำหรับ response
	result := make([]NotiItemWithNames, len(rows))
	for i, r := range rows {
		result[i] = NotiItemWithNames{
			ID:              r.ID,
			PatientID:       r.PatientID,
			MyMedicineID:    r.MyMedicineID,
			GroupID:         r.GroupID,
			GroupName:       r.GroupName,
			NotiInfoID:      r.NotiInfoID,
			NotifyDate:      r.NotifyDate,
			NotifyTime:      r.NotifyTime,
			TakenStatus:     r.TakenStatus,
			NotifyStatus:    r.NotifyStatus,
			MedName:         r.MedName,
			AmountPerTime:   r.AmountPerTime,
			FormID:          r.FormID,
			FormName:        r.FormName,
			UnitID:          r.UnitID,
			UnitName:        r.UnitName,
			InstructionID:   r.InstructionID,
			InstructionName: r.InstructionName,
			HasSymptom:      false,
		}
	}

	return result, nil
}

// ===================================================================
//                Update: เปลี่ยนสถานะ taken/notify
// ===================================================================

// MarkNotiItemTaken: เซ็ต/ยกเลิกสถานะ “ทานแล้ว” (รองรับอัปเดตทั้งกลุ่มใน slot)
func MarkNotiItemTaken(db *gorm.DB, patientID, notiItemID uint, taken bool) (*models.NotiItem, error) {
	var item models.NotiItem
	// จำกัดสิทธิ์: ต้องเป็นของผู้ป่วยคนนี้เท่านั้น
	if err := db.Where("id = ? AND patient_id = ?", notiItemID, patientID).First(&item).Error; err != nil {
		return nil, err
	}

	nowLocal := time.Now().In(time.Local)
	var takenTime *time.Time
	if taken {
		takenTime = &nowLocal
	} else {
		takenTime = nil
	}

	// ถ้าเป็นรายการใน "กลุ่ม" -> อัปเดตทั้งชุดใน slot เดียวกัน
	if item.GroupID != nil {
		if err := db.Model(&models.NotiItem{}).
			Where("patient_id = ? AND group_id = ? AND noti_info_id = ? AND notify_date = ? AND notify_time = ?",
				patientID, *item.GroupID, item.NotiInfoID, item.NotifyDate, item.NotifyTime).
			Updates(map[string]any{
				"taken_status":  taken,
				"taken_time_at": takenTime,
			}).Error; err != nil {
			return nil, err
		}
	} else {
		// เดี่ยว: อัปเดตเฉพาะรายการนี้
		if err := db.Model(&item).Updates(map[string]any{
			"taken_status":  taken,
			"taken_time_at": takenTime,
		}).Error; err != nil {
			return nil, err
		}
	}

	// อ่านกลับรายการที่กด (ใช้ตอบกลับ)
	if err := db.Where("id = ? AND patient_id = ?", notiItemID, patientID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// MarkNotiItemNotified: เซ็ต/ยกเลิก “แจ้งเตือนแล้ว” (รองรับอัปเดตทั้งกลุ่มใน slot)
func MarkNotiItemNotified(db *gorm.DB, patientID, notiItemID uint, notified bool) (*models.NotiItem, error) {
	var item models.NotiItem
	// จำกัดสิทธิ์: ต้องเป็นของผู้ป่วยคนนี้เท่านั้น
	if err := db.Where("id = ? AND patient_id = ?", notiItemID, patientID).First(&item).Error; err != nil {
		return nil, err
	}

	if item.GroupID != nil {
		// อัปเดตทั้งชุดใน slot เดียวกัน
		if err := db.Model(&models.NotiItem{}).
			Where("patient_id = ? AND group_id = ? AND noti_info_id = ? AND notify_date = ? AND notify_time = ?",
				patientID, *item.GroupID, item.NotiInfoID, item.NotifyDate, item.NotifyTime).
			Update("notify_status", notified).Error; err != nil {
			return nil, err
		}
	} else {
		// เดี่ยว
		if err := db.Model(&item).Update("notify_status", notified).Error; err != nil {
			return nil, err
		}
	}

	// อ่านกลับรายการที่กด
	if err := db.Where("id = ? AND patient_id = ?", notiItemID, patientID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ListNotiFormats(db *gorm.DB) ([]models.NotiFormat, error) {
	var formats []models.NotiFormat
	if err := db.Find(&formats).Error; err != nil {
		return nil, err
	}
	return formats, nil
}
