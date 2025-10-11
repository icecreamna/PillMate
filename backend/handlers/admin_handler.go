package handlers

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/fouradithep/pillmate/models"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// --- helpers ---

func signJWT(claims jwt.MapClaims) (string, error) {
	secret := os.Getenv("jwtSecretKey")
	if secret == "" {
		return "", errors.New("jwt secret not set")
	}
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = claims
	return t.SignedString([]byte(secret))
}

// --- login แบบปกติ: ดึงจาก DB เท่านั้น ---

// LoginAdmin ทำการล็อกอินจากตาราง web_admins เท่านั้น (ไม่มีตรวจ .env)
// - เทียบ username แบบ case-insensitive
// - เทียบรหัสผ่านกับ hash ใน DB ด้วย bcrypt
// - ออก JWT พร้อม role จาก DB และ admin_id
func LoginAdmin(db *gorm.DB, username, password string) (string, *models.WebAdmin, error) {
	u := strings.TrimSpace(username)
	p := strings.TrimSpace(password)

	if u == "" || p == "" {
		return "", nil, errors.New("username/password required")
	}

	var admin models.WebAdmin
	// เทียบแบบ case-insensitive: LOWER(username) = LOWER(?)
	if err := db.Where("LOWER(username) = LOWER(?)", u).First(&admin).Error; err != nil {
		// ไม่บอกว่าไม่พบผู้ใช้ เพื่อความปลอดภัย -> รวมเป็น invalid credentials
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(p)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"aud":      "admin-app",
		"role":     admin.Role, // เช่น "superadmin" / "doctor" / "staff"
		"admin_id": admin.ID,
		"iat":      now.Unix(),
		"exp":      now.Add(72 * time.Hour).Unix(), // ปรับอายุโทเค็นได้ตามต้องการ
	}
	// ถ้าต้องการความเข้ากันได้เดิมกับฝั่ง UI ที่คาดหวัง doctor_id
	if strings.EqualFold(admin.Role, "doctor") {
		claims["doctor_id"] = admin.ID
	}

	tok, err := signJWT(claims)
	if err != nil {
		return "", nil, err
	}

	// เคลียร์ password ก่อนคืนค่า
	admin.Password = ""
	return tok, &admin, nil
}
