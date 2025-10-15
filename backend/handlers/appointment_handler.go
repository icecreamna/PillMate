package handlers

import (
	"errors"
	"gorm.io/gorm"

	// "github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/models"
	"time"
)

// // ใบนัดล่าสุดของผู้ใช้ (เวลาถูกฟอร์แมต "YYYY-MM-DD"/"HH:MM" แล้วใน DTO Mobile)
// func MobileGetLatestAppointment(db *gorm.DB, patientID uint) (*dto.MobileAppointmentResponse, error) {
// 	if patientID == 0 {
// 		return nil, errors.New("invalid patient")
// 	}

// 	var p models.Patient
// 	if err := db.Select("id_card_number").First(&p, patientID).Error; err != nil {
// 		return nil, err
// 	}
// 	if p.IDCardNumber == "" {
// 		return nil, errors.New("missing id_card_number for this account")
// 	}

// 	var appt models.Appointment
// 	if err := db.
// 		Where("id_card_number = ?", p.IDCardNumber).
// 		Order("appointment_date DESC, appointment_time DESC, id DESC").
// 		First(&appt).Error; err != nil {
// 		return nil, err
// 	}

// 	resp := dto.ToMobileAppointmentResponse(appt)
// 	return &resp, nil
// }

// // อ่านใบนัดตาม id (เวลาถูกฟอร์แมต "YYYY-MM-DD"/"HH:MM" แล้วใน DTO Mobile)
// func MobileGetAppointmentByID(db *gorm.DB, patientID uint, appointmentID uint) (*dto.MobileAppointmentResponse, error) {
// 	if patientID == 0 || appointmentID == 0 {
// 		return nil, errors.New("invalid arguments")
// 	}

// 	var p models.Patient
// 	if err := db.Select("id_card_number").First(&p, patientID).Error; err != nil {
// 		return nil, err
// 	}
// 	if p.IDCardNumber == "" {
// 		return nil, errors.New("missing id_card_number for this account")
// 	}

// 	var appt models.Appointment
// 	if err := db.
// 		Where("id = ? AND id_card_number = ?", appointmentID, p.IDCardNumber).
// 		First(&appt).Error; err != nil {
// 		return nil, err
// 	}

// 	resp := dto.ToMobileAppointmentResponse(appt)
// 	return &resp, nil
// }

func GetNextAppointment(db *gorm.DB, patientID uint) (*models.Appointment, error) {
	// ใช้โซนเวลาไทยให้สม่ำเสมอ
	appLoc, _ := time.LoadLocation("Asia/Bangkok")

	// เวลา "ตอนนี้" ตามไทย
	nowLocal := time.Now().In(appLoc)

	// "วันนี้" เที่ยงคืน UTC ของวันนั้น (ให้ตรงกับวิธีเก็บคอลัมน์ DATE)
	todayUTC := time.Date(
		nowLocal.Year(), nowLocal.Month(), nowLocal.Day(),
		0, 0, 0, 0, time.UTC,
	)

	// เวลา HH:MM:SS ของ "ตอนนี้" (ยึดวันที่ 2000-01-01 @Bangkok แล้ว .UTC()) ให้ชนิดตรงกับคอลัมน์ TIME
	nowHMUTC := time.Date(
		2000, 1, 1,
		nowLocal.Hour(), nowLocal.Minute(), nowLocal.Second(), 0,
		appLoc,
	).UTC()

	var nextAppointment models.Appointment

	// ✅ Query นัดที่ยังไม่ถึงวัน/เวลา
	err := db.
		Where(`
			patient_id = ?
			AND (
				appointment_date > ?
				OR (appointment_date = ? AND appointment_time > ?)
			)
		`, patientID, todayUTC, todayUTC, nowHMUTC).
		Order("appointment_date ASC, appointment_time ASC, id ASC").
		First(&nextAppointment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &nextAppointment, nil
}
