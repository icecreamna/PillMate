package handlers

import (
	"errors"
	"time"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

// READ: ใบนัด 1 รายการ โดยยืนยันสิทธิ์ด้วย (id, id_card_number)
func GetAppointment(db *gorm.DB, appointmentID uint, idCardNumber string) (*models.Appointment, error) {
	var appt models.Appointment
	if err := db.
		Preload("WebAdmin").
		Where("id = ? AND id_card_number = ?", appointmentID, idCardNumber).
		First(&appt).Error; err != nil {
		return nil, err
	}
	return &appt, nil
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
