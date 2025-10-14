package handlers

import (
	"errors"
	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/models"
)

// ใบนัดล่าสุดของผู้ใช้ (เวลาถูกฟอร์แมต "YYYY-MM-DD"/"HH:MM" แล้วใน DTO Mobile)
func MobileGetLatestAppointment(db *gorm.DB, patientID uint) (*dto.MobileAppointmentResponse, error) {
	if patientID == 0 {
		return nil, errors.New("invalid patient")
	}

	var p models.Patient
	if err := db.Select("id_card_number").First(&p, patientID).Error; err != nil {
		return nil, err
	}
	if p.IDCardNumber == "" {
		return nil, errors.New("missing id_card_number for this account")
	}

	var appt models.Appointment
	if err := db.
		Where("id_card_number = ?", p.IDCardNumber).
		Order("appointment_date DESC, appointment_time DESC, id DESC").
		First(&appt).Error; err != nil {
		return nil, err
	}

	resp := dto.ToMobileAppointmentResponse(appt)
	return &resp, nil
}

// อ่านใบนัดตาม id (เวลาถูกฟอร์แมต "YYYY-MM-DD"/"HH:MM" แล้วใน DTO Mobile)
func MobileGetAppointmentByID(db *gorm.DB, patientID uint, appointmentID uint) (*dto.MobileAppointmentResponse, error) {
	if patientID == 0 || appointmentID == 0 {
		return nil, errors.New("invalid arguments")
	}

	var p models.Patient
	if err := db.Select("id_card_number").First(&p, patientID).Error; err != nil {
		return nil, err
	}
	if p.IDCardNumber == "" {
		return nil, errors.New("missing id_card_number for this account")
	}

	var appt models.Appointment
	if err := db.
		Where("id = ? AND id_card_number = ?", appointmentID, p.IDCardNumber).
		First(&appt).Error; err != nil {
		return nil, err
	}

	resp := dto.ToMobileAppointmentResponse(appt)
	return &resp, nil
}

func GetNextAppointment(db *gorm.DB, patientID uint) (*models.Appointment, error) {
	now := time.Now().In(time.Local)
	today := now.Format("2006-01-02")

	var nextAppointment models.Appointment

	// ✅ Query นัดที่ยังไม่ถึงวัน/เวลา
	err := db.Select(`
		appointment_date,
		('2000-01-01 ' || TO_CHAR(appointment_time, 'HH24:MI:SS'))::timestamp AS appointment_time,
		note
	`).
		Where(`
			patient_id = ?
			AND (
				appointment_date > ?
				OR (appointment_date = ? AND appointment_time > ?)
			)
		`, patientID, today, today, now.Format("15:04:05")).
		Order("appointment_date ASC, appointment_time ASC").
		First(&nextAppointment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &nextAppointment, nil
}
