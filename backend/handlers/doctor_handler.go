package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateDoctor สร้างบัญชีหมอใหม่ (role=doctor)
func CreateDoctor(db *gorm.DB, in *dto.CreateDoctorDTO, actorID uint) (*models.WebAdmin, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}
	u := strings.TrimSpace(in.Username)
	p := strings.TrimSpace(in.Password)
	fn := strings.TrimSpace(in.FirstName)
	ln := strings.TrimSpace(in.LastName)
	if u == "" || p == "" {
		return nil, errors.New("username and password are required")
	}

	// ตรวจซ้ำ username แบบเร็ว
	var cnt int64
	if err := db.Model(&models.WebAdmin{}).Where("username = ?", u).Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, errors.New("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	doc := &models.WebAdmin{
		Username:  u,
		Password:  string(hashed),
		FirstName: fn,
		LastName:  ln,
		Role:      "doctor",
		CreatedAt: now,
		UpdatedAt: now,
		// ถ้ามีฟิลด์ CreatedByID/UpdatedByID ในโมเดลค่อยใส่เพิ่ม
	}
	if err := db.Create(doc).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

// ListDoctors ค้นหา/แบ่งหน้า รายชื่อหมอ
func ListDoctors(db *gorm.DB, q string, page, pageSize int) ([]models.WebAdmin, int64, error) {
	var out []models.WebAdmin
	var total int64

	tx := db.Model(&models.WebAdmin{}).Where("role = ?", "doctor")

	if s := strings.TrimSpace(q); s != "" {
		like := "%" + s + "%"
		tx = tx.Where("username ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?", like, like, like)
	}

	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := tx.Order("id DESC").Limit(pageSize).Offset(offset).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

// GetDoctorByID ดึงหมอจาก id
func GetDoctorByID(db *gorm.DB, id uint) (*models.WebAdmin, error) {
	var doc models.WebAdmin
	if err := db.Where("role = ?", "doctor").First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

// UpdateDoctor อัปเดตข้อมูลหมอ
func UpdateDoctor(db *gorm.DB, id uint, in *dto.UpdateDoctorDTO) (*models.WebAdmin, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	var doc models.WebAdmin
	if err := db.Where("role = ?", "doctor").First(&doc, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]any{
		"updated_at": time.Now(),
	}

	if in.Username != nil {
		u := strings.TrimSpace(*in.Username)
		if u == "" {
			return nil, errors.New("username cannot be empty")
		}
		// ตรวจซ้ำ username (ยกเว้นกรณีเดิม)
		var cnt int64
		if err := db.Model(&models.WebAdmin{}).
			Where("username = ? AND id <> ?", u, id).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errors.New("username already exists")
		}
		updates["username"] = u
	}

	if in.FirstName != nil {
		updates["first_name"] = strings.TrimSpace(*in.FirstName)
	}
	if in.LastName != nil {
		updates["last_name"] = strings.TrimSpace(*in.LastName)
	}
	if in.Password != nil {
		pw := strings.TrimSpace(*in.Password)
		if pw != "" {
			hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}
			updates["password"] = string(hashed)
		}
	}

	if err := db.Model(&doc).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := db.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

// DeleteDoctor ลบหมอ (soft delete ถ้าโมเดลมี gorm.DeletedAt)
func DeleteDoctor(db *gorm.DB, id uint) error {
	// จำกัดเฉพาะ role=doctor กันเผลอลบ superadmin
	return db.Where("role = ?", "doctor").Delete(&models.WebAdmin{}, id).Error
}

// ResetDoctorPassword ตั้งรหัสผ่านใหม่ (สำหรับแอดมินรีเซ็ต)
func ResetDoctorPassword(db *gorm.DB, id uint, newPassword string) (*models.WebAdmin, error) {
	pw := strings.TrimSpace(newPassword)
	if pw == "" {
		return nil, errors.New("password cannot be empty")
	}

	var doc models.WebAdmin
	if err := db.Where("role = ?", "doctor").First(&doc, id).Error; err != nil {
		return nil, err
	}

	// กันการตั้งเป็นรหัสเดิมเป๊ะ ๆ
	if err := bcrypt.CompareHashAndPassword([]byte(doc.Password), []byte(pw)); err == nil {
		return nil, errors.New("new password must be different from current password")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if err := db.Model(&doc).Updates(map[string]any{
		"password":   string(hashed),
		"updated_at": time.Now(),
	}).Error; err != nil {
		return nil, err
	}

	if err := db.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

