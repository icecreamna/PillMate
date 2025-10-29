package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

// ใช้สำหรับ RESPONSE เท่านั้น (ตัด field ที่ไม่จำเป็นออก)
type MedicineInfoDTO struct {
	ID            uint      `json:"id"`
	MedName       string    `json:"med_name"`
	GenericName   string    `json:"generic_name"`
	Properties    string    `json:"properties"`
	Strength      string    `json:"strength"`
	FormID        uint      `json:"form_id"`
	UnitID        *uint     `json:"unit_id,omitempty"`
	InstructionID *uint     `json:"instruction_id,omitempty"`
	MedStatus     string    `json:"med_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// แปลง model -> DTO (ตัวเดียว)
func ToMedicineInfoDTO(m *models.MedicineInfo) MedicineInfoDTO {
	return MedicineInfoDTO{
		ID:            m.ID,
		MedName:       m.MedName,
		GenericName:   m.GenericName,
		Properties:    m.Properties,
		Strength:      m.Strength,
		FormID:        m.FormID,
		UnitID:        m.UnitID,
		InstructionID: m.InstructionID,
		MedStatus:     m.MedStatus,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

// แปลง list ของ model -> DTO list
func ToMedicineInfoDTOs(list []models.MedicineInfo) []MedicineInfoDTO {
	out := make([]MedicineInfoDTO, 0, len(list))
	for i := range list {
		out = append(out, ToMedicineInfoDTO(&list[i]))
	}
	return out
}
