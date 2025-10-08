package dto

import (
	"time"
	
	"github.com/fouradithep/pillmate/models"
)

type AppointmentDTO struct {
	ID               uint   `json:"id"`
	IDCardNumber     string `json:"id_card_number"`
	AppointmentDate  string `json:"appointment_date"` // "YYYY-MM-DD"
	AppointmentTime  string `json:"appointment_time"` // "HH:MM"

	HospitalID       uint   `json:"hospital_id"`
	HospitalName     string `json:"hospital_name,omitempty"`

	DoctorID         uint   `json:"doctor_id"`
	DoctorFirstName  string `json:"doctor_first_name,omitempty"`
	DoctorLastName   string `json:"doctor_last_name,omitempty"`

	Note             string `json:"note,omitempty"`
}

func AppointmentToDTO(a models.Appointment) AppointmentDTO {
	return AppointmentDTO{
		ID:               a.ID,
		IDCardNumber:     a.IDCardNumber,
		AppointmentDate:  a.AppointmentDate.Format("2006-01-02"),
		AppointmentTime:  a.AppointmentTime.In(time.UTC).Format("15:04"),

		HospitalID:       a.HospitalID,
		HospitalName:     a.Hospital.HospitalName,

		DoctorID:         a.DoctorID,
		DoctorFirstName:  a.WebAdmin.FirstName,
		DoctorLastName:   a.WebAdmin.LastName,

		Note:             a.Note,
	}
}

func AppointmentsToDTO(list []models.Appointment) []AppointmentDTO {
	out := make([]AppointmentDTO, 0, len(list))
	for _, ap := range list {
		out = append(out, AppointmentToDTO(ap))
	}
	return out
}
