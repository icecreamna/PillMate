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

func ListNotiItems(db *gorm.DB, patientID uint, f ListNotiItemsFilter) ([]models.NotiItem, error) {
	q := db.Model(&models.NotiItem{}).
		Preload("Patient").
		Preload("MyMedicine").
		Preload("Group").
		Preload("NotiInfo").
		Preload("Form").
		Preload("Unit").
		Preload("Instruction")

	// บังคับกรองตาม patient_id เสมอ (เมิน f.PatientID เพื่อกัน override)
	q = q.Where("patient_id = ?", patientID)

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

	var items []models.NotiItem
	if err := q.Order("notify_date, notify_time").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
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