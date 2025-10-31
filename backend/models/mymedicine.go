package models

import(
	"gorm.io/gorm"
	"time"
)

// กลุ่ม 
type Group struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PatientID uint           `gorm:"not null;index" json:"patient_id"`
	GroupName string         `gorm:"type:varchar(255);not null" json:"group_name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
// ยาของฉัน
type MyMedicine struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	PatientID       uint           `gorm:"not null;index" json:"patient_id"`

	// อ้างอิงข้อมูลยา (optional เพื่อให้รองรับทั้ง manual/hospital)
	MedicineInfoID  *uint          `gorm:"index" json:"medicine_info_id,omitempty"`

	MedName         string         `gorm:"type:varchar(255);not null" json:"med_name"`
	Properties      string         `gorm:"not null" json:"properties"`

	FormID          uint           `gorm:"not null;index" json:"form_id"`             // รูปแบบยา (บังคับมี)
	UnitID          *uint          `gorm:"index" json:"unit_id,omitempty"`            // nullable
	InstructionID   *uint          `gorm:"index" json:"instruction_id,omitempty"`     // nullable

	AmountPerTime   string         `gorm:"not null" json:"amount_per_time"`
	TimesPerDay     string         `gorm:"not null" json:"times_per_day"`

	// แหล่งที่มา: manual หรือ hospital
	Source          string         `gorm:"type:varchar(16);check:source IN ('manual','hospital');not null" json:"source"`

	// อ้างอิงใบสั่งยา: รองรับโมเดลหัว+รายการ
	PrescriptionID      *uint      `gorm:"index" json:"prescription_id,omitempty"`
	PrescriptionItemID  *uint      `gorm:"index" json:"prescription_item_id,omitempty"`

	// จัดกลุ่ม (เช่น ในแอปรวมเป็นกลุ่ม)
	GroupID         *uint          `gorm:"index" json:"group_id,omitempty"`

	// ==== เพิ่มฟิลด์วันหมดอายุ ====
	ExpireAt        *time.Time     `gorm:"index" json:"expire_at,omitempty"`

	// Relations
	Group        Group        `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Patient      Patient      `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Form         Form         `gorm:"foreignKey:FormID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	Unit         Unit         `gorm:"foreignKey:UnitID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	Instruction  Instruction  `gorm:"foreignKey:InstructionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`
	MedicineInfo MedicineInfo `gorm:"foreignKey:MedicineInfoID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"-"`

	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
