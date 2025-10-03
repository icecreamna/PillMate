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

// MarkNotiItemTaken: เซ็ต/ยกเลิกสถานะ “ทานแล้ว”
func MarkNotiItemTaken(db *gorm.DB, patientID, notiItemID uint, taken bool) (*models.NotiItem, error) {
	var item models.NotiItem
	// จำกัดสิทธิ์: ต้องเป็นรายการของผู้ป่วยคนนี้เท่านั้น
	if err := db.Where("id = ? AND patient_id = ?", notiItemID, patientID).First(&item).Error; err != nil {
		return nil, err
	}

	updateFields := map[string]any{
		"taken_status": taken,
	}
	if taken {
		nowLocal := time.Now().In(time.Local)
		updateFields["taken_time_at"] = &nowLocal
	} else {
		// ถ้าอยากล้างเวลา เมื่อยกเลิก ให้ตั้งค่าเป็น NULL ได้เพราะเป็น *time.Time
		updateFields["taken_time_at"] = nil
	}

	if err := db.Model(&item).Updates(updateFields).Error; err != nil {
		return nil, err
	}
	// อ่านกลับยืนยันผล โดยยังคงจำกัดสิทธิ์เดิม
	if err := db.Where("id = ? AND patient_id = ?", notiItemID, patientID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// MarkNotiItemNotified: เซ็ต/ยกเลิก “แจ้งเตือนแล้ว”
func MarkNotiItemNotified(db *gorm.DB, patientID, notiItemID uint, notified bool) (*models.NotiItem, error) {
	var item models.NotiItem
	// จำกัดสิทธิ์: ต้องเป็นรายการของผู้ป่วยคนนี้เท่านั้น
	if err := db.Where("id = ? AND patient_id = ?", notiItemID, patientID).First(&item).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&item).Update("notify_status", notified).Error; err != nil {
		return nil, err
	}
	// อ่านกลับยืนยันผล โดยยังคงจำกัดสิทธิ์เดิม
	if err := db.Where("id = ? AND patient_id = ?", notiItemID, patientID).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
