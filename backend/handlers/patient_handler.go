package handlers

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/fouradithep/pillmate/mailer"
	"github.com/fouradithep/pillmate/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// -------- Patients: register/login --------

func CreatePatient(db *gorm.DB, patient *models.Patient) error {
	// 1) hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(patient.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	patient.Password = string(hashedPassword)

	// 2) defaults
	if strings.TrimSpace(patient.VerificationStatus) == "" {
		patient.VerificationStatus = "unverified"
	}
	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()

	// 3) save
	if err := db.Create(patient).Error; err != nil {
		return err
	}
	return nil
}

func LoginPatient(db *gorm.DB, patient *models.Patient) (string, error) {
	// get user by email
	var selected models.Patient
	if err := db.Where("email = ?", patient.Email).First(&selected).Error; err != nil {
		return "", err
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(selected.Password), []byte(patient.Password)); err != nil {
		return "", err
	}

	// issue JWT
	jwtSecretKey := os.Getenv("jwtSecretKey") // set in .env
	if jwtSecretKey == "" {
		return "", errors.New("JWT secret key not set")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["patient_id"] = selected.ID
	// claims["exp"] = time.Now().Add(72 * time.Hour).Unix() // จะมี token จนกว่าจะ logout

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}
	return t, nil
}

// -------- OTP helpers --------

// สุ่ม OTP 6 หลัก
func GenerateOTP6() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// ออก OTP + บันทึก + ส่งอีเมล
func IssueOTP(db *gorm.DB, patientID uint, ttl time.Duration) (*models.VerificationCode, error) {
	// หาอีเมลผู้ใช้
	var p models.Patient
	if err := db.Select("id", "email").First(&p, patientID).Error; err != nil {
		return nil, err
	}
	if strings.TrimSpace(p.Email) == "" {
		return nil, errors.New("patient has no email")
	}

	// (ออปชัน) ยกเลิก OTP เดิมที่ยังไม่หมดอายุ
	if err := RevokeActiveOTP(db, patientID); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// ไม่ถึงกับ fatal แต่ควร log ภายนอกถ้าต้องการ
	}

	// สร้าง OTP
	raw, err := GenerateOTP6()
	if err != nil {
		return nil, err
	}

	// บันทึก DB
	vc := &models.VerificationCode{
		OTPCode:   raw,                         // เก็บตามโมเดลเดิมของลูก (plaintext)
		PatientID: patientID,
		ExpiresAt: time.Now().Add(ttl),
	}
	if err := db.Create(vc).Error; err != nil {
		return nil, err
	}

	// ส่งอีเมล
	m, err := mailer.New()
	if err != nil {
		return nil, err
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "PillMate"
	}
	subject := "[" + appName + "] รหัสยืนยัน OTP ของคุณ"
	html := fmt.Sprintf(`
		<div style="font-family:system-ui,-apple-system,Segoe UI,Roboto,Helvetica,Arial">
			<h2>%s</h2>
			<p>รหัสยืนยันของคุณคือ:</p>
			<div style="font-size:28px;font-weight:700;letter-spacing:4px">%s</div>
			<p>รหัสมีอายุใช้งาน %d นาที</p>
			<p>หากคุณไม่ได้ร้องขอ สามารถเพิกเฉยอีเมลนี้ได้</p>
		</div>
	`, appName, raw, int(ttl.Minutes()))
	text := fmt.Sprintf("รหัส OTP ของคุณคือ %s (หมดอายุใน %d นาที)", raw, int(ttl.Minutes()))

	if err := m.Send(p.Email, subject, html, text); err != nil {
		return nil, err
	}

	return vc, nil
}

// ยกเลิกโค้ดเก่าที่ยังไม่หมดอายุทั้งหมดก่อนออกใหม่ (soft delete)
func RevokeActiveOTP(db *gorm.DB, patientID uint) error {
	return db.Where("patient_id = ? AND expires_at > ?", patientID, time.Now()).
		Delete(&models.VerificationCode{}).Error
}

// ตรวจสอบรหัส OTP ล่าสุดของผู้ใช้ (ใช้โมเดลเดิม: ลบ soft-delete หลังใช้)
func VerifyOTP(db *gorm.DB, patientID uint, input string) error {
	var vc models.VerificationCode

	// ดึง OTP ล่าสุดที่ยังไม่ถูกลบ
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

	// เทียบแบบ constant-time
	if subtle.ConstantTimeCompare([]byte(vc.OTPCode), []byte(input)) != 1 {
		return errors.New("invalid otp")
	}

	// ลบ OTP (soft delete) เพื่อกันใช้ซ้ำ
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

// รีเซ็ตรหัสผ่าน
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
