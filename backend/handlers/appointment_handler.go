package handlers

import (
	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

// READ: ดึงใบนัด "รายการเดียว" โดยยืนยันด้วย appointmentID + id_card_number
func GetAppointment(db *gorm.DB, appointmentID uint, idCardNumber string) (*models.Appointment, error) {
	var appt models.Appointment
	if err := db.
		Where("id = ? AND id_card_number = ?", appointmentID, idCardNumber).
		Preload("Hospital").
		Preload("WebAdmin").
		First(&appt).Error; err != nil {
		return nil, err
	}
	return &appt, nil
}
