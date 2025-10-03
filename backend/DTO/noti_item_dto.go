package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

type NotiItemDTO struct {
	ID            uint   `json:"id"`
	PatientID     uint   `json:"patient_id"`
	MyMedicineID  uint   `json:"my_medicine_id"`
	GroupID       *uint  `json:"group_id,omitempty"`   // เพิ่มไว้ใช้รวมการ์ด
	NotiInfoID    uint   `json:"noti_info_id"`
	NotifyDate    string `json:"notify_date"`          // "YYYY-MM-DD"
	NotifyTime    string `json:"notify_time"`          // "HH:MM"
	TakenStatus   bool   `json:"taken_status"`
	NotifyStatus  bool   `json:"notify_status"`
	MedName       string `json:"med_name"`
	GroupName     string `json:"group_name,omitempty"`
	AmountPerTime string `json:"amount_per_time"`
	FormID        uint   `json:"form_id"`
	UnitID        *uint  `json:"unit_id,omitempty"`
	InstructionID *uint  `json:"instruction_id,omitempty"`
}

func NotiItemToDTO(m models.NotiItem) NotiItemDTO {
	return NotiItemDTO{
		ID:            m.ID,
		PatientID:     m.PatientID,
		MyMedicineID:  m.MyMedicineID,
		GroupID:       m.GroupID,                                  // map ออกมา
		NotiInfoID:    m.NotiInfoID,
		NotifyDate:    m.NotifyDate.Format("2006-01-02"),
		NotifyTime:    m.NotifyTime.In(time.UTC).Format("15:04"),  // HH:MM คงที่
		TakenStatus:   m.TakenStatus,
		NotifyStatus:  m.NotifyStatus,
		MedName:       m.MedName,
		GroupName:     m.GroupName,
		AmountPerTime: m.AmountPerTime,
		FormID:        m.FormID,
		UnitID:        m.UnitID,
		InstructionID: m.InstructionID,
	}
}

func NotiItemsToDTO(items []models.NotiItem) []NotiItemDTO {
	out := make([]NotiItemDTO, 0, len(items))
	for _, it := range items {
		out = append(out, NotiItemToDTO(it))
	}
	return out
}
