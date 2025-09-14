package handlers

import(
	"time"

	"github.com/fouradithep/pillmate/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	 "os"
	// "errors"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"math/big"
	

)

func CreatePatient(db *gorm.DB, patient *models.Patient) error {
	// 1. เข้ารหัสรหัสผ่านก่อนบันทึก
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(patient.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	patient.Password = string(hashedPassword)

	// 2. ตั้งค่า default values ถ้าอยากกำหนดเอง
	if patient.VerificationStatus == "" {
		patient.VerificationStatus = "unverified"
	}
	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()

	// 3. บันทึกลง DB
	result := db.Create(patient)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func LoginPatient(db *gorm.DB, patient *models.Patient) (string, error) {
	// get user from email
	selectedPatient := new(models.Patient)
	result := db.Where("email = ?", patient.Email).First(selectedPatient)

	if result.Error != nil {
		return "", result.Error
	}

	// compare passwprd
	err := bcrypt.CompareHashAndPassword([]byte(selectedPatient.Password), []byte(patient.Password))

	if err != nil {
		return "", err
	}

	// pass = return jwt
	jwtSecretKey := os.Getenv("jwtSecretKey") //in .env
	
	if jwtSecretKey == "" {
		return "", errors.New("JWT secret key not set")
	
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["patient_id"] = selectedPatient.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

// สุ่ม OTP 6 หลัก
func GenerateOTP6() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// ออก OTP แล้ว "บันทึกไว้ก่อน" เฉย ๆ (ยังไม่ต้องส่งให้ผู้ใช้)
func IssueOTP(db *gorm.DB, patientID uint, ttl time.Duration) (*models.VerificationCode, error) {
	raw, err := GenerateOTP6()
	if err != nil {
		return nil, err
	}
	vc := &models.VerificationCode{
		OTPCode:   raw,
		PatientID: patientID,
		ExpiresAt: time.Now().Add(ttl),
	}
	if err := db.Create(vc).Error; err != nil {
		return nil, err
	}
	return vc, nil
}

// ยกเลิกโค้ดเก่าที่ยังไม่หมดอายุทั้งหมดก่อนออกใหม่
func RevokeActiveOTP(db *gorm.DB, patientID uint) error {
	return db.Where("patient_id = ? AND expires_at > ?", patientID, time.Now()).
		Delete(&models.VerificationCode{}).Error
}

// VerifyOTP ตรวจสอบรหัส OTP ล่าสุดของผู้ใช้
func VerifyOTP(db *gorm.DB, patientID uint, input string) error {
	var vc models.VerificationCode

	// ดึง OTP ล่าสุดของผู้ใช้ที่ยังไม่ถูกลบ
	if err := db.Where("patient_id = ? AND deleted_at IS NULL", patientID).
		Order("id DESC").
		First(&vc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("no otp found, please request a new one")
		}
		return err
	}

	// ตรวจสอบหมดอายุ
	if time.Now().After(vc.ExpiresAt) {
		return errors.New("otp expired")
	}

	// ตรวจสอบว่า OTP ถูกต้อง (ใช้ ConstantTimeCompare กัน timing attack)
	if subtle.ConstantTimeCompare([]byte(vc.OTPCode), []byte(input)) != 1 {
		return errors.New("invalid otp")
	}

	// ลบ OTP ทิ้ง (soft delete) เพื่อกันใช้ซ้ำ
	if err := db.Delete(&vc).Error; err != nil {
		return err
	}

	// อัปเดตสถานะผู้ใช้เป็น verified
	if err := db.Model(&models.Patient{}).
		Where("id = ?", patientID).
		Update("verification_status", "verified").Error; err != nil {
		return err
	}

	return nil
}

func UpdatePatientPassword(db *gorm.DB, patientID uint, newPlain string) error {
    hashed, err := bcrypt.GenerateFromPassword([]byte(newPlain), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    return db.Model(&models.Patient{}).
        Where("id = ?", patientID).
        Updates(map[string]interface{}{
            "password":   string(hashed),
            "updated_at": time.Now(),
        }).Error
}