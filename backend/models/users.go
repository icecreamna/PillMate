package models

import(
	"gorm.io/gorm"
	"time"
)
type Patient struct {
	ID 		 		uint `gorm:"primaryKey" json:"id"`
	Email    		string `gorm:"unique" json:"email"`
	Password 		string `json:"password"`
	IDCardNumber 	*string `gorm:"type:char(13);unique" json:"id_card_number"`
	FirstName 		string `gorm:"type:varchar(255)" json:"first_name"`
    LastName  		string `gorm:"type:varchar(255)" json:"last_name"`
	PhoneNumber 	*string `gorm:"type:varchar(10);unique" json:"phone_number"`
	VerificationStatus string `gorm:"check:verification_status IN ('unverified','verified');default:'unverified'" json:"verification_status"`
	CreatedAt 		time.Time      `json:"created_at"`
	UpdatedAt 		time.Time      `json:"updated_at"`
	DeletedAt 		gorm.DeletedAt `gorm:"index" json:"-"`
}


type WebAdmin struct {
	ID 		 	uint `gorm:"primaryKey" json:"id"`
	Username    string `gorm:"unique" json:"username"`
	Password 	string `json:"password"`
	FirstName 	string `gorm:"type:varchar(255)" json:"first_name"`
    LastName  	string `gorm:"type:varchar(255)" json:"last_name"`
	Role      	string `gorm:"type:varchar(50);default:'doctor'" json:"role"` // เช่น "superadmin", "doctor" 
	CreatedAt 	time.Time      `json:"created_at"`
	UpdatedAt 	time.Time      `json:"updated_at"`
	DeletedAt 	gorm.DeletedAt `gorm:"index" json:"-"`		
}

type HospitalPatient struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	IDCardNumber  string         `gorm:"type:varchar(13);not null;uniqueIndex:uniq_active_idcard,where:deleted_at IS NULL" json:"id_card_number"`
	PatientCode   string         `gorm:"type:varchar(20);not null;uniqueIndex:uniq_patient_code,where:deleted_at IS NULL" json:"patient_code"` // รหัสผู้ป่วยภายในโรงพยาบาล
	FirstName     string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName      string         `gorm:"type:varchar(100);not null" json:"last_name"`
	PhoneNumber   string         `gorm:"type:varchar(10);not null;uniqueIndex:uniq_active_phone,where:deleted_at IS NULL" json:"phone_number"`
	BirthDay      time.Time      `gorm:"type:date;not null" json:"birth_day"`
	Gender        string         `gorm:"type:varchar(10);not null" json:"gender"` // ค่า: ชาย, หญิง (กำหนด CHECK ด้านล่าง)
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}