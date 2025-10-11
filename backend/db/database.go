package db

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/joho/godotenv"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	
)

var DB *gorm.DB

func Init() {
	
	// โหลดค่า .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("ไม่สามารถโหลดไฟล์ .env:", err)
	}

	// อ่านค่าจาก environment
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		host, user, password, dbname, port)
		fmt.Println("DSN:", dsn)

	// ตั้ง timezone ของ Go เป็นไทย
	loc, _ := time.LoadLocation("Asia/Bangkok")
	time.Local = loc

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal("ไม่สามารถเชื่อมต่อฐานข้อมูล:", err)
	}

	// สร้างตารางอัตโนมัติ
	err = db.AutoMigrate(
		&models.Patient{}, 
		&models.VerificationCode{},
		&models.Form{},
		&models.Unit{},
		&models.FormUnit{},
		&models.Instruction{},
		// &models.Hospital{},
		&models.MedicineInfo{},
		&models.Prescription{},
		
		&models.Appointment{},
		// &models.AppointmentNoti{},
		&models.NotiLog{},
		&models.Group{},
		&models.MyMedicine{},
		&models.NotiFormat{},
		&models.NotiInfo{},
		&models.NotiItem{},
		&models.Symptom{},
	)
	if err != nil {
		log.Fatal("AutoMigrate ล้มเหลว:", err)
	}

	// เพิ่มข้อมูลจากseed
	SeedInitialData(db)

	DB = db
	fmt.Println("เชื่อมต่อฐานข้อมูลสำเร็จด้วย GORM")

	// Test timezone
	var med models.MedicineInfo
	if err := db.First(&med).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Println("CreatedAt Go Local Time:", med.CreatedAt)

}
