package handlers

import (
	"errors"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

// -------------------------
// Response DTO (ไม่ผูก gorm)
// -------------------------
type PatientProfileResponse struct {
	ID                 uint    `json:"id"`
	Email              string  `json:"email"`
	IDCardNumber       *string `json:"id_card_number,omitempty"`
	FirstName          string  `json:"first_name"`
	LastName           string  `json:"last_name"`
	PhoneNumber        *string `json:"phone_number,omitempty"`
	VerificationStatus string  `json:"verification_status"`
	CreatedAt          any     `json:"created_at"`
	UpdatedAt          any     `json:"updated_at"`
	// patient_code จากตาราง hospital_patients (null ถ้าไม่มี)
	PatientCode *string `json:"patient_code"`
}

// แปลง models.Patient -> PatientProfileResponse พร้อมเติม patient_code
func toPatientProfileResponse(db *gorm.DB, p models.Patient) (PatientProfileResponse, error) {
	code, err := findPatientCodeByIDCard(db, p.IDCardNumber)
	if err != nil {
		return PatientProfileResponse{}, err
	}
	return PatientProfileResponse{
		ID:                 p.ID,
		Email:              p.Email,
		IDCardNumber:       p.IDCardNumber,
		FirstName:          p.FirstName,
		LastName:           p.LastName,
		PhoneNumber:        p.PhoneNumber,
		VerificationStatus: p.VerificationStatus,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
		PatientCode:        code, // nil ถ้าไม่พบ
	}, nil
}

// ค้นหา patient_code จากตาราง hospital_patients ด้วย id_card_number (respect soft delete)
// คืนค่าเป็น *string (nil = ไม่พบ/ไม่มี)
func findPatientCodeByIDCard(db *gorm.DB, idCardPtr *string) (*string, error) {
	if idCardPtr == nil || *idCardPtr == "" {
		return nil, nil
	}
	var hp struct {
		PatientCode string
	}
	if err := db.
		Table("hospital_patients").
		Select("patient_code").
		Where("id_card_number = ? AND deleted_at IS NULL", *idCardPtr).
		Take(&hp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &hp.PatientCode, nil
}

// READ: โปรไฟล์ของตัวเอง (+ patient_code)
func GetPatient(db *gorm.DB, patientID uint) (*PatientProfileResponse, error) {
	var patient models.Patient
	if err := db.Where("id = ?", patientID).First(&patient).Error; err != nil {
		return nil, err
	}
	// ไม่ส่งรหัสผ่านกลับ
	patient.Password = ""

	resp, err := toPatientProfileResponse(db, patient)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UPDATE: อัปเดตได้เฉพาะ 4 ฟิลด์ที่อนุญาต (+ คืน patient_code)
func UpdatePatientBasic(db *gorm.DB, patientID uint, payload *models.Patient) (*PatientProfileResponse, error) {
	// รวบรวมเฉพาะฟิลด์ที่อนุญาตให้อัปเดต
	updateFields := map[string]any{}

	// ป้องกันการเขียนค่าเป็นค่าว่างโดยไม่ตั้งใจ
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
		var current models.Patient
		if err := db.Where("id = ?", patientID).First(&current).Error; err != nil {
			return nil, err
		}
		current.Password = ""
		resp, err := toPatientProfileResponse(db, current)
		if err != nil {
			return nil, err
		}
		return &resp, nil
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
	var updated models.Patient
	if err := db.Where("id = ?", patientID).First(&updated).Error; err != nil {
		return nil, err
	}
	updated.Password = ""

	resp, err := toPatientProfileResponse(db, updated)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
