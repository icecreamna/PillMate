package db

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func getenv(k, def string) string {
	if v := strings.TrimSpace(os.Getenv(k)); v != "" {
		return v
	}
	return def
}

func mustEnv(keys ...string) {
	for _, k := range keys {
		if strings.TrimSpace(os.Getenv(k)) == "" {
			log.Fatalf("missing required env: %s", k)
		}
	}
}

func Init() {
	// ----- โหลดไฟล์ env หลายไฟล์ตามโหมด -----
	// ลองโหลด .env (ถ้ามี) + .env.common + .env.<mode>
	_ = godotenv.Load(".env") // optional
	mode := getenv("APP_MODE", "all")
	_ = godotenv.Load(".env.common")
	_ = godotenv.Load(".env." + mode) // .env.admin / .env.mobile / .env.all (แล้วแต่คุณสร้าง)

	// ----- อ่านค่าจำเป็นและตรวจ -----
	mustEnv("DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := getenv("DB_SSLMODE", "disable")

	// ตั้ง timezone ของ Go เป็นไทย
	if loc, err := time.LoadLocation("Asia/Bangkok"); err == nil {
		time.Local = loc
	}

	// ----- DSN (ไม่พิมพ์ password ออก log) -----
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		host, user, password, dbname, port, sslmode,
	)
	safeDSN := fmt.Sprintf(
		"host=%s user=%s password=**** dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		host, user, dbname, port, sslmode,
	)
	fmt.Println("DSN:", safeDSN)

	// ----- เปิด GORM -----
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal("ไม่สามารถเชื่อมต่อฐานข้อมูล:", err)
	}

	// ----- Health check DB -----
	if err := db.Exec("SELECT 1").Error; err != nil {
		log.Fatal("DB ping failed:", err)
	}

	// ----- AutoMigrate (ให้รันเฉพาะฝั่งที่กำหนด) -----
	// ตั้ง RUN_MIGRATIONS=true ใน .env.admin เท่านั้น
	if strings.EqualFold(getenv("RUN_MIGRATIONS", "false"), "true") {
		if err := db.AutoMigrate(
			&models.Patient{},
			&models.VerificationCode{},
			&models.Form{},
			&models.Unit{},
			&models.FormUnit{},
			&models.Instruction{},
			&models.Hospital{},
			&models.MedicineInfo{},
			&models.Prescription{},
			&models.PrescriptionItem{},
			&models.Appointment{},
			// &models.AppointmentNoti{},
			// &models.NotiLog{},
			&models.Group{},
			&models.MyMedicine{},
			&models.NotiFormat{},
			&models.NotiInfo{},
			&models.NotiItem{},
			&models.Symptom{},

			// ฝั่งเว็บ ถ้ามี เช่น:
			&models.WebAdmin{},
			&models.HospitalPatient{},
		); err != nil {
			log.Fatal("AutoMigrate ล้มเหลว:", err)
		}
		// seed (ถ้าต้อง) — ขยับมาอยู่หลัง migrate และทำเฉพาะฝั่งที่รัน migrate
		SeedInitialData(db)
	}

	DB = db
	fmt.Println("เชื่อมต่อฐานข้อมูลสำเร็จด้วย GORM")

	// ----- แสดงเวลา NOW() จาก DB (ไม่ fatal ถ้าว่าง) -----
	var now time.Time
	if err := db.Raw("SELECT NOW()").Scan(&now).Error; err == nil {
		fmt.Println("DB NOW():", now.In(time.Local))
	}
}
