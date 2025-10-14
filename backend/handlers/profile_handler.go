package handlers

import (
	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

// READ: โปรไฟล์ของตัวเอง
func GetPatient(db *gorm.DB, patientID uint) (*models.Patient, error) {
	var patient models.Patient
	if err := db.Where("id = ?", patientID).First(&patient).Error; err != nil {
		return nil, err
	}
	// ไม่ส่งรหัสผ่านกลับ
	patient.Password = ""
	return &patient, nil
}

// UPDATE: อัปเดตได้เฉพาะ 4 ฟิลด์ที่อนุญาต
func UpdatePatientBasic(db *gorm.DB, patientID uint, payload *models.Patient) (*models.Patient, error) {
	// รวบรวมเฉพาะฟิลด์ที่อนุญาตให้อัปเดต
	updateFields := map[string]any{}

	// หมายเหตุ: ป้องกันการเขียนค่าเป็นค่าว่างโดยไม่ตั้งใจ
	if payload.IDCardNumber != nil {
		updateFields["id_card_number"] = payload.IDCardNumber
	}
	if payload.FirstName != "" {
		updateFields["first_name"] = payload.FirstName
	}
	if payload.LastName != "" {
		updateFields["last_name"] = payload.LastName
	}
	if payload.PhoneNumber != nil {
		updateFields["phone_number"] = payload.PhoneNumber
	}

	// ถ้าไม่มีอะไรให้อัปเดต ให้คืนค่าปัจจุบัน
	if len(updateFields) == 0 {
		var currentPatient models.Patient
		if err := db.Where("id = ?", patientID).First(&currentPatient).Error; err != nil {
			return nil, err
		}
		currentPatient.Password = ""
		return &currentPatient, nil
	}

	// อัปเดตเฉพาะเรคคอร์ดของตัวเอง
	result := db.Model(&models.Patient{}).
		Where("id = ?", patientID).
		Updates(updateFields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// โหลดกลับและล้าง password ออกจาก response
	var updatedPatient models.Patient
	if err := db.Where("id = ?", patientID).First(&updatedPatient).Error; err != nil {
		return nil, err
	}
	updatedPatient.Password = ""
	return &updatedPatient, nil
}
