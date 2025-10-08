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
	ID 		 		uint `gorm:"primaryKey" json:"id"`
	IDCardNumber 	string `gorm:"type:char(13);unique" json:"id_card_number"`
	FirstName 		string `gorm:"type:varchar(255)" json:"first_name"`
    LastName  		string `gorm:"type:varchar(255)" json:"last_name"`
	PhoneNumber 	string `gorm:"type:varchar(10);unique" json:"phone_number"`
	BirthDay        time.Time `gorm:"type:date" json:"birth_day"`
	Age             int `gorm:"not null" json:"age"` //อายุ
	Gender			string `gorm:"not null" json:"gender"` //ชาย, หญิง
	CreatedAt 		time.Time      `json:"created_at"`
	UpdatedAt 		time.Time      `json:"updated_at"`
	DeletedAt 		gorm.DeletedAt `gorm:"index" json:"-"`

}