package handlers

import (
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
