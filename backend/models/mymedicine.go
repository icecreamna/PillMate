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
    ID             uint      `gorm:"primaryKey" json:"id"`
    PatientID      uint      `gorm:"not null" json:"patient_id"`
    MedName        string    `gorm:"type:varchar(255);not null" json:"med_name"`
    Properties     string    `gorm:"not null" json:"properties"`
    FormID         uint      `gorm:"not null" json:"form_id"` // รูปแบบยา (บังคับมี)
    UnitID         *uint     `json:"unit_id,omitempty"`       // ← เปลี่ยนเป็น pointer เพื่อให้เป็น NULL ได้
    InstructionID  *uint     `json:"instruction_id,omitempty"`// ← เปลี่ยนเป็น pointer เพื่อให้เป็น NULL ได้
    AmountPerTime  string    `gorm:"not null" json:"amount_per_time"`
    TimesPerDay    string    `gorm:"not null" json:"times_per_day"`
    Source         string    `gorm:"check:source IN ('manual','hospital')" json:"source"`
    PrescriptionID *uint     `gorm:"index" json:"prescription_id,omitempty"`
    GroupID        *uint     `gorm:"index" json:"group_id,omitempty"`

    Group          Group     `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
    Patient        Patient   `gorm:"foreignKey:PatientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Form           Form      `gorm:"foreignKey:FormID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
    Unit           Unit      `gorm:"foreignKey:UnitID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
    Instruction    Instruction `gorm:"foreignKey:InstructionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

    CreatedAt      time.Time      `json:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
    DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
