package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

// ========== Request DTOs ==========

// รายการยาแต่ละตัวในใบสั่ง (สำหรับสร้าง/แทนที่)
type CreatePrescriptionItemDTO struct {
	MedicineInfoID uint   `json:"medicine_info_id"`          // ต้องมีอยู่จริง
	AmountPerTime  string `json:"amount_per_time"`           // เช่น "1 เม็ด"
	TimesPerDay    string `json:"times_per_day"`             // เช่น "2 ครั้ง"
}

// สร้างใบสั่งยา (หัว) พร้อมรายการยา 1..N รายการ
type CreatePrescriptionDTO struct {
	IDCardNumber  string                      `json:"id_card_number"`                 // ต้องเป็นตัวเลข 13 หลัก (ตรวจใน handler)
	DoctorID      uint                        `json:"doctor_id"`                      // ต้องมีอยู่จริง และ role=doctor (ไม่ส่งได้ถ้าจะอิงจาก token ฝั่ง route)
	SyncUntil     *time.Time                  `json:"sync_until,omitempty"`           // optional; ไม่ส่ง = +60 วัน (ตั้งใน hook)
	AppSyncStatus *bool                       `json:"app_sync_status,omitempty"`      // optional; default=false
	Items         []CreatePrescriptionItemDTO `json:"items"`                          // ต้องมี >= 1
}

// อัปเดตใบสั่งยา (หัวเอกสาร) — ไม่รวมการแก้ไขรายการยา (items)
// หมายเหตุ: ถ้าต้องการแก้ไข items แนะนำทำ endpoint แยก (เช่น replace ทั้งชุด หรือ add/update/delete รายการทีละตัว)
type UpdatePrescriptionDTO struct {
	IDCardNumber  *string    `json:"id_card_number,omitempty"`
	DoctorID      *uint      `json:"doctor_id,omitempty"`
	SyncUntil     *time.Time `json:"sync_until,omitempty"`
	AppSyncStatus *bool      `json:"app_sync_status,omitempty"`
}

// (ตัวเลือกเพิ่มเติม หากภายหลังต้องการแก้ไขรายการยาเป็นชุดเดียว)
// type ReplacePrescriptionItemsDTO struct {
// 	Items []CreatePrescriptionItemDTO `json:"items"` // แทนที่ทั้งชุด
// }

// ========== Response DTOs ==========

type PrescriptionItemResponse struct {
	ID             uint      `json:"id"`
	MedicineInfoID uint      `json:"medicine_info_id"`
	AmountPerTime  string    `json:"amount_per_time"`
	TimesPerDay    string    `json:"times_per_day"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
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
			CreatedAt:      it.CreatedAt,
			UpdatedAt:      it.UpdatedAt,
		})
	}
	return out
}
