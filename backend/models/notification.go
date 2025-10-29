package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

// รูปแบบการแจ้งเตือน
type NotiFormat struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	FormatName string `gorm:"type:varchar(255);not null;unique" json:"format_name"` //ชื่อรูปแบบ

}

// ข้อมูลแจ้งเตือน
type NotiInfo struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	PatientID    uint      `gorm:"not null;index" json:"patient_id"` 
	MyMedicineID *uint     `json:"my_medicine_id"`
	GroupID      *uint     `json:"group_id"`
	StartDate    time.Time `gorm:"type:date" json:"start_date"`
	EndDate      time.Time `gorm:"type:date" json:"end_date"`
	NotiFormatID uint      `gorm:"not null" json:"noti_format_id"`

	Times         *pq.StringArray `gorm:"type:text[]" json:"times"`        // ["08:00","12:00","20:00"]
	IntervalHours *int            `json:"interval_hours"`                  // แจ้งเตือนทุกกี่ชั่วโมง
	TimesPerDay   *int            `json:"times_per_day"`                   //กี่ครั้งต่อวัน
	IntervalDay   *int            `json:"interval_day"`                    // แจ้งเตือนทุกกี่วัน
	CyclePattern  *pq.Int64Array  `gorm:"type:int[]" json:"cycle_pattern"` // ตัวอย่าง: [21,7]

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Patient    Patient    `gorm:"foreignKey:PatientID"` 
	MyMedicine MyMedicine `gorm:"foreignKey:MyMedicineID"`
	Group      Group      `gorm:"foreignKey:GroupID"`
	NotiFormat NotiFormat `gorm:"foreignKey:NotiFormatID"`
}

// รายการแจ้งเตือน
type NotiItem struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	PatientID uint `gorm:"not null;index:idx_patient_date" json:"patient_id"`

	// ---- ชุดคีย์ป้องกันซ้ำ (ต้องมีค่าเสมอ) ----
	MyMedicineID uint      `gorm:"not null;uniqueIndex:noti_items_uniq_active,where:deleted_at IS NULL" json:"my_medicine_id"`
	NotifyDate   time.Time `gorm:"type:date;not null;uniqueIndex:noti_items_uniq_active,where:deleted_at IS NULL;index:idx_patient_date" json:"notify_date"`
	NotifyTime   time.Time `gorm:"type:time;not null;uniqueIndex:noti_items_uniq_active,where:deleted_at IS NULL" json:"notify_time"`
	NotiInfoID   uint      `gorm:"not null;uniqueIndex:noti_items_uniq_active,where:deleted_at IS NULL" json:"noti_info_id"`

	// ---- ตัวเลือก ----
	GroupID       *uint  `json:"group_id"` // nullable
	MedName       string `gorm:"type:varchar(255);not null" json:"med_name"`
	GroupName     string `gorm:"type:varchar(255);not null" json:"group_name"`
	AmountPerTime string `gorm:"not null" json:"amount_per_time"`
	FormID        uint   `gorm:"not null" json:"form_id"`
	UnitID        *uint  `json:"unit_id"`        // nullable
	InstructionID *uint  `json:"instruction_id"` // nullable

	TakenStatus  bool       `gorm:"not null;default:false;index" json:"taken_status"`
	TakenTimeAt  *time.Time `json:"taken_time_at"`
	NotifyStatus bool       `gorm:"not null;default:false;index" json:"notify_status"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Patient     Patient     `gorm:"foreignKey:PatientID"`
	MyMedicine  MyMedicine  `gorm:"foreignKey:MyMedicineID"`
	Group       Group       `gorm:"foreignKey:GroupID"`
	NotiInfo    NotiInfo    `gorm:"foreignKey:NotiInfoID"`
	Form        Form        `gorm:"foreignKey:FormID"`
	Unit        Unit        `gorm:"foreignKey:UnitID"`
	Instruction Instruction `gorm:"foreignKey:InstructionID"`
}
