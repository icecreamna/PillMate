package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

type NotiItemDTO struct {
	ID            uint   `json:"id"`
	PatientID     uint   `json:"patient_id"`
	MyMedicineID  uint   `json:"my_medicine_id"`
	NotiInfoID    uint   `json:"noti_info_id"`
	NotifyDate    string `json:"notify_date"`
	NotifyTime    string `json:"notify_time"`
	TakenStatus   bool   `json:"taken_status"`
	NotifyStatus  bool   `json:"notify_status"`
	MedName       string `json:"med_name"`
	GroupName     string `json:"group_name,omitempty"`
	AmountPerTime string `json:"amount_per_time"`
	FormID        uint   `json:"form_id"`
	UnitID        *uint  `json:"unit_id,omitempty"`
	InstructionID *uint  `json:"instruction_id,omitempty"`

	
	HasSymptom bool  `json:"has_symptom"`
	SymptomID  *uint `json:"symptom_id,omitempty"`
}

func NotiItemToDTO(m models.NotiItem) NotiItemDTO {
	return NotiItemDTO{
		ID:            m.ID,
		PatientID:     m.PatientID,
		MyMedicineID:  m.MyMedicineID,
		NotiInfoID:    m.NotiInfoID,
		NotifyDate:    m.NotifyDate.Format("2006-01-02"),
		NotifyTime:    m.NotifyTime.In(time.UTC).Format("15:04"),
		TakenStatus:   m.TakenStatus,
		NotifyStatus:  m.NotifyStatus,
		MedName:       m.MedName,
		GroupName:     m.GroupName,
		AmountPerTime: m.AmountPerTime,
		FormID:        m.FormID,
		UnitID:        m.UnitID,
		InstructionID: m.InstructionID,
		HasSymptom:    false,
		SymptomID:     nil,
	}
}

func NotiItemsToDTO(items []models.NotiItem) []NotiItemDTO {
	out := make([]NotiItemDTO, 0, len(items))
	for _, it := range items {
		out = append(out, NotiItemToDTO(it))
	}
	return out
}

// แปลงพร้อมข้อมูลอาการ (symMap: noti_item_id -> symptom_id)
func NotiItemsToDTOWithSymptoms(items []models.NotiItem, symMap map[uint]uint) []NotiItemDTO {
	out := make([]NotiItemDTO, 0, len(items))
	for _, it := range items {
		dto := NotiItemToDTO(it)
		if sid, ok := symMap[it.ID]; ok {
			dto.HasSymptom = true
			dto.SymptomID = &sid
		}
		out = append(out, dto)
	}
	return out
}
   