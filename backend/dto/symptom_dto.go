package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

type SymptomDTO struct {
	ID           uint   `json:"id"`
	PatientID    uint   `json:"patient_id"`
	MyMedicineID uint   `json:"my_medicine_id"`
	GroupID      *uint  `json:"group_id,omitempty"` 
	NotiItemID   uint   `json:"noti_item_id"`
	SymptomNote  string `json:"symptom_note"`

	CreatedAt string `json:"created_at"` 
	UpdatedAt string `json:"updated_at"`
}

func SymptomToDTO(s models.Symptom) SymptomDTO {
	return SymptomDTO{
		ID:           s.ID,
		PatientID:    s.PatientID,
		MyMedicineID: s.MyMedicineID,
		GroupID:      s.GroupID, // ถ้าโมเดลคุณเป็น uint ให้ใช้: &tmp := s.GroupID; แล้วใส่ tmp
		NotiItemID:   s.NotiItemID,
		SymptomNote:  s.SymptomNote,
		CreatedAt:    s.CreatedAt.In(time.Local).Format(time.RFC3339),
		UpdatedAt:    s.UpdatedAt.In(time.Local).Format(time.RFC3339),
	}
}

func SymptomsToDTO(list []models.Symptom) []SymptomDTO {
	out := make([]SymptomDTO, 0, len(list))
	for _, s := range list {
		out = append(out, SymptomToDTO(s))
	}
	return out
}
