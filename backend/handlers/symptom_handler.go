package handlers

import (
	"time"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

// ===============================
//            CREATE
// ===============================

// CreateSymptom: บันทึกโน้ตอาการใหม่ (จำกัดสิทธิ์ด้วย patientID)
// - ผูกกับ NotiItem ของผู้ป่วยคนนี้เท่านั้น
// - เติม MyMedicineID / GroupID ให้สอดคล้องกับ NotiItem ที่อ้างถึง
func CreateSymptom(db *gorm.DB, patientID uint, inputSymptom *models.Symptom) (*models.Symptom, error) {
	// กัน owner/primary key
	inputSymptom.ID = 0
	inputSymptom.PatientID = patientID

	// โหลด NotiItem ที่อ้างถึง และยืนยันว่าเป็นของผู้ป่วยนี้
	var linkedNotiItem models.NotiItem
	if err := db.Where("id = ? AND patient_id = ?", inputSymptom.NotiItemID, patientID).
		First(&linkedNotiItem).Error; err != nil {
		return nil, err
	}

	// เติม/ตรวจ MyMedicineID ให้ตรงกับ NotiItem
	if inputSymptom.MyMedicineID == 0 {
		inputSymptom.MyMedicineID = linkedNotiItem.MyMedicineID
	} else if inputSymptom.MyMedicineID != linkedNotiItem.MyMedicineID {
		return nil, gorm.ErrInvalidData
	}

	// เติม/ตรวจ GroupID ให้ตรงกับ NotiItem (ถ้ามี)
	if inputSymptom.GroupID == nil {
		if linkedNotiItem.GroupID != nil {
			inputSymptom.GroupID = linkedNotiItem.GroupID
		}
	} else {
		if linkedNotiItem.GroupID == nil || *inputSymptom.GroupID != *linkedNotiItem.GroupID {
			return nil, gorm.ErrInvalidData
		}
	}

	// ต้องมีโน้ตอาการ
	if inputSymptom.SymptomNote == "" {
		return nil, gorm.ErrInvalidData
	}

	if err := db.Create(inputSymptom).Error; err != nil {
		return nil, err
	}
	return inputSymptom, nil
}

// ===============================
//         READ (ONE / LIST)
// ===============================

// GetSymptom: อ่านอาการรายการเดียวของผู้ป่วย
func GetSymptom(db *gorm.DB, patientID, symptomID uint) (*models.Symptom, error) {
	var symptom models.Symptom
	if err := db.
		Where("id = ? AND patient_id = ?", symptomID, patientID).
		First(&symptom).Error; err != nil {
		return nil, err
	}
	return &symptom, nil
}

// ฟิลเตอร์สำหรับ ListSymptoms
type ListSymptomsFilter struct {
	MyMedicineID *uint
	GroupID      *uint
	NotiItemID   *uint
	CreatedFrom  *string // "YYYY-MM-DD"
	CreatedTo    *string // "YYYY-MM-DD"
}

// ListSymptoms: ดึงอาการทั้งหมดของผู้ป่วย พร้อมฟิลเตอร์พื้นฐาน
func ListSymptoms(db *gorm.DB, patientID uint, filter ListSymptomsFilter) ([]models.Symptom, error) {
	queryBuilder := db.Model(&models.Symptom{}).
		Where("patient_id = ?", patientID)

	if filter.MyMedicineID != nil {
		queryBuilder = queryBuilder.Where("my_medicine_id = ?", *filter.MyMedicineID)
	}
	if filter.GroupID != nil {
		queryBuilder = queryBuilder.Where("group_id = ?", *filter.GroupID)
	}
	if filter.NotiItemID != nil {
		queryBuilder = queryBuilder.Where("noti_item_id = ?", *filter.NotiItemID)
	}

	const ymd = "2006-01-02"
	if filter.CreatedFrom != nil && *filter.CreatedFrom != "" {
		if fromTime, err := time.ParseInLocation(ymd, *filter.CreatedFrom, time.Local); err == nil {
			queryBuilder = queryBuilder.Where("created_at >= ?", fromTime)
		}
	}
	if filter.CreatedTo != nil && *filter.CreatedTo != "" {
		if toTime, err := time.ParseInLocation(ymd, *filter.CreatedTo, time.Local); err == nil {
			dayEnd := time.Date(toTime.Year(), toTime.Month(), toTime.Day(), 23, 59, 59, 0, toTime.Location())
			queryBuilder = queryBuilder.Where("created_at <= ?", dayEnd)
		}
	}

	var symptomsList []models.Symptom
	if err := queryBuilder.Order("created_at DESC, id DESC").Find(&symptomsList).Error; err != nil {
		return nil, err
	}
	return symptomsList, nil
}

// ===============================
//            UPDATE
// ===============================

// UpdateSymptom: อัปเดตเฉพาะโน้ตอาการ (SymptomNote)
func UpdateSymptom(db *gorm.DB, patientID, symptomID uint, inputSymptom *models.Symptom) (*models.Symptom, error) {
	var existingSymptom models.Symptom
	if err := db.Where("id = ? AND patient_id = ?", symptomID, patientID).
		First(&existingSymptom).Error; err != nil {
		return nil, err
	}

	updatedFields := map[string]any{}
	if inputSymptom.SymptomNote != "" {
		updatedFields["symptom_note"] = inputSymptom.SymptomNote
	}

	if len(updatedFields) == 0 {
		return &existingSymptom, nil
	}

	if err := db.Model(&existingSymptom).Updates(updatedFields).Error; err != nil {
		return nil, err
	}

	if err := db.Where("id = ? AND patient_id = ?", symptomID, patientID).
		First(&existingSymptom).Error; err != nil {
		return nil, err
	}
	return &existingSymptom, nil
}

// ===============================
//            DELETE
// ===============================

// DeleteSymptom: ลบอาการของผู้ป่วย (soft delete ถ้ามี DeletedAt)
func DeleteSymptom(db *gorm.DB, patientID, symptomID uint) error {
	deleteResult := db.Where("id = ? AND patient_id = ?", symptomID, patientID).
		Delete(&models.Symptom{})
	if deleteResult.Error != nil {
		return deleteResult.Error
	}
	if deleteResult.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
