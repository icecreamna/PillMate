package models

import(
	"gorm.io/gorm"
	"time"
)

// รูปแบบยา
type Form struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    FormName string `gorm:"type:varchar(100);not null;unique" json:"form_name"` // เช่น เม็ด, แคปซูล, ยาน้ำ, ยาฉีด

	Units []Unit `gorm:"many2many:form_units;" json:"units"` // Many-to-Many
}

// หน่วยยา
type Unit struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    UnitName string `gorm:"type:varchar(100);not null;unique" json:"unit_name"` // เช่น เม็ด, แคปซูล, ช้อนชา, cc

	Forms []Form `gorm:"many2many:form_units;" json:"forms"` // Many-to-Many
}

// เชื่อมโยงรูปแบบยากับหน่วยยา
type FormUnit struct {
    FormID uint `gorm:"primaryKey"`
    UnitID uint `gorm:"primaryKey"`
}

// ช่วงเวลาใช้ยา
type Instruction struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    InstructionName string `gorm:"type:varchar(100);not null;unique" json:"instruction_name"` // เช่น ก่อนอาหาร, หลังอาหาร
}

// ข้อมูลยา
type MedicineInfo struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	MedName string `gorm:"type:varchar(255);not null;unique" json:"med_name"` //ชื่อทางการค้า
	GenericName string `gorm:"type:varchar(255);not null;unique" json:"generic_name"` //ชื่อสามัญ
	Properties string `gorm:"not null" json:"properties"` //รายละเอียด สรรพคุณ
	Strength  string `gorm:"type:varchar(255);not null" json:"strength"` //ความแรง

	FormID uint `gorm:"not null" json:"form_id"`
    UnitID uint `gorm:"default:null" json:"unit_id"`
	InstructionID uint `gorm:"default:null" json:"instruction_id"`
	

    Form Form `gorm:"foreignKey:FormID"`
    Unit Unit `gorm:"foreignKey:UnitID"`
	Instruction Instruction `gorm:"foreignKey:InstructionID"`
	

	MedStatus string `gorm:"check:med_status IN ('active','inactive');default:'active'" json:"med_status"` //สถานะว่ายายังมีการใช้อยู่ไหม

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
