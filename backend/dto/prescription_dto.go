package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

// ========== Request DTOs ==========

// รายการยาแต่ละตัวในใบสั่ง (สำหรับสร้าง/แทนที่)
type CreatePrescriptionItemDTO struct {
	MedicineInfoID uint    `json:"medicine_info_id"`          // ต้องมีอยู่จริง
	AmountPerTime  string  `json:"amount_per_time"`           // เช่น "1 เม็ด"
	TimesPerDay    string  `json:"times_per_day"`             // เช่น "2 ครั้ง"
	// ฟิลด์ใหม่ (optional) — เก็บเป็น "YYYY-MM-DD"
	StartDate      *string `json:"start_date,omitempty"`
	EndDate        *string `json:"end_date,omitempty"`
	// ไม่รับ expire_date จาก client (คำนวณใน hook)
	Note           *string `json:"note,omitempty"`
}

// สร้างใบสั่งยา (หัว) พร้อมรายการยา 1..N รายการ
type CreatePrescriptionDTO struct {
	IDCardNumber  string                      `json:"id_card_number"`
	DoctorID      uint                        `json:"doctor_id"`
	SyncUntil     *time.Time                  `json:"sync_until,omitempty"`
	AppSyncStatus *bool                       `json:"app_sync_status,omitempty"`
	Items         []CreatePrescriptionItemDTO `json:"items"` // ต้องมี >= 1
}

// อัปเดตใบสั่งยา (หัวเอกสาร)
type UpdatePrescriptionDTO struct {
	IDCardNumber  *string    `json:"id_card_number,omitempty"`
	DoctorID      *uint      `json:"doctor_id,omitempty"`
	SyncUntil     *time.Time `json:"sync_until,omitempty"`
	AppSyncStatus *bool      `json:"app_sync_status,omitempty"`
}

// ========== Response DTOs ==========

type PrescriptionItemResponse struct {
	ID             uint      `json:"id"`
	MedicineInfoID uint      `json:"medicine_info_id"`
	AmountPerTime  string    `json:"amount_per_time"`
	TimesPerDay    string    `json:"times_per_day"`
	// ฟิลด์ใหม่ (optional) — "YYYY-MM-DD"
	StartDate   *string   `json:"start_date,omitempty"`
	EndDate     *string   `json:"end_date,omitempty"`
	ExpireDate  *string   `json:"expire_date,omitempty"` // คำนวณจาก EndDate + 1 วัน (hook)
	Note        *string   `json:"note,omitempty"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PrescriptionResponse struct {
	ID            uint                       `json:"id"`
	IDCardNumber  string                     `json:"id_card_number"`
	DoctorID      uint                       `json:"doctor_id"`
	AppSyncStatus bool                       `json:"app_sync_status"`
	SyncUntil     time.Time                  `json:"sync_until"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	Items         []PrescriptionItemResponse `json:"items"`
}

func NewPrescriptionResponse(m models.Prescription) PrescriptionResponse {
	out := PrescriptionResponse{
		ID:            m.ID,
		IDCardNumber:  m.IDCardNumber,
		DoctorID:      m.DoctorID,
		AppSyncStatus: m.AppSyncStatus,
		SyncUntil:     m.SyncUntil,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		Items:         make([]PrescriptionItemResponse, 0, len(m.Items)),
	}
	for _, it := range m.Items {
		out.Items = append(out.Items, PrescriptionItemResponse{
			ID:             it.ID,
			MedicineInfoID: it.MedicineInfoID,
			AmountPerTime:  it.AmountPerTime,
			TimesPerDay:    it.TimesPerDay,

			StartDate:      it.StartDate,
			EndDate:        it.EndDate,
			ExpireDate:     it.ExpireDate, // ได้มาจาก hook ของโมเดล
			Note:           it.Note,

			CreatedAt:      it.CreatedAt,
			UpdatedAt:      it.UpdatedAt,
		})
	}
	return out
}
