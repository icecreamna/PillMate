package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

type AppointmentDTO struct {
	ID             uint   `json:"id"`
	IDCardNumber   string `json:"id_card_number"`
	AppointmentDate string `json:"appointment_date"` // "YYYY-MM-DD"
	AppointmentTime string `json:"appointment_time"` // "HH:MM"
	HospitalID     uint   `json:"hospital_id"`
	DoctorID       uint   `json:"doctor_id"`
	Note           string `json:"note"`
}

func AppointmentToDTO(a models.Appointment) AppointmentDTO {
	return AppointmentDTO{
		ID:              a.ID,
		IDCardNumber:    a.IDCardNumber,
		AppointmentDate: a.AppointmentDate.In(time.Local).Format("2006-01-02"),
		// เวลาเก็บเป็น time-only แนะนำฟอร์แมตด้วย UTC เพื่อเลี่ยง offset แปลก (+06:42 ในปี 0001)
		AppointmentTime: a.AppointmentTime.In(time.UTC).Format("15:04"),
		HospitalID:      a.HospitalID,
		DoctorID:        a.DoctorID,
		Note:            a.Note,
	}
}

func AppointmentsToDTO(list []models.Appointment) []AppointmentDTO {
	out := make([]AppointmentDTO, 0, len(list))
	for _, a := range list {
		out = append(out, AppointmentToDTO(a))
	}
	return out
}
