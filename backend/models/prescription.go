package models

import (
	"time"

	"gorm.io/gorm"
)

// ใบสั่งยา (หัวเอกสาร) — 1 ใบสั่งยา สามารถมีหลายรายการยา (PrescriptionItem)
type Prescription struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	IDCardNumber string         `gorm:"type:varchar(13);not null;index" json:"id_card_number"`
	DoctorID     uint           `gorm:"not null;index" json:"doctor_id"`

	// timestamps + soft delete
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	AppSyncStatus bool      `gorm:"default:false" json:"app_sync_status"` // false=ยังไม่ซิงค์
	SyncUntil     time.Time `gorm:"not null;index" json:"sync_until"`     // ซิงก์ได้ถึงวันที่กำหนด (ไม่ส่ง = +60 วัน)

	// รายการยาในใบสั่งนี้ (ลบหัว → ลบลูกด้วย; hard delete จะ CASCADE โดย DB)
	Items []PrescriptionItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`

	// ความสัมพันธ์อื่น (เลือก preload ตามต้องการ)
	WebAdmin WebAdmin `gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

// รายการยาในใบสั่ง (หนึ่งแถวต่อยาหนึ่งตัว/หนึ่งโดส)
type PrescriptionItem struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	PrescriptionID uint           `gorm:"not null;index" json:"prescription_id"`
	MedicineInfoID uint           `gorm:"not null;index" json:"medicine_info_id"`
	AmountPerTime  string         `gorm:"not null" json:"amount_per_time"` // เช่น "1 เม็ด"
	TimesPerDay    string         `gorm:"not null" json:"times_per_day"`   // เช่น "2 ครั้ง"

	// timestamps + soft delete
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// ความสัมพันธ์อ้างอิงข้อมูลยา
	MedicineInfo MedicineInfo `gorm:"foreignKey:MedicineInfoID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

// ตั้งค่า SyncUntil อัตโนมัติถ้าไม่กำหนดมา
func (p *Prescription) BeforeCreate(tx *gorm.DB) error {
	if p.SyncUntil.IsZero() {
		p.SyncUntil = time.Now().AddDate(0, 0, 60)
	}
	return nil
}