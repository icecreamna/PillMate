package handlers

import (
	"errors"
	"strings"
	"unicode"

	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/models"
)

// =====================
// Helpers
// =====================

func onlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func norm(s string) string { return strings.TrimSpace(s) }

func validGender(g string) bool {
	g = norm(g)
	return g == "ชาย" || g == "หญิง"
}

// func now() time.Time { return time.Now() }

// =====================
// CREATE
// =====================

func CreateHospitalPatient(db *gorm.DB, in *dto.CreateHospitalPatientDTO) (*dto.HospitalPatientResponse, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	idcard := onlyDigits(in.IDCardNumber)
	phone := onlyDigits(in.PhoneNumber)
	fn := norm(in.FirstName)
	ln := norm(in.LastName)
	g := norm(in.Gender)

	if idcard == "" || phone == "" || fn == "" || ln == "" || in.BirthDay.IsZero() || g == "" {
		return nil, errors.New("missing required fields")
	}
	if len(idcard) != 13 {
		return nil, errors.New("id_card_number must be 13 digits")
	}
	if len(phone) != 10 {
		return nil, errors.New("phone_number must be 10 digits")
	}
	if !validGender(g) {
		return nil, errors.New(`gender must be "ชาย" or "หญิง"`)
	}

	// ตรวจซ้ำแบบรวดเร็ว (respect soft delete: deleted_at IS NULL)
	var count int64
	if err := db.Model(&models.HospitalPatient{}).
		Where("id_card_number = ? AND deleted_at IS NULL", idcard).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("id_card_number already exists")
	}
	if err := db.Model(&models.HospitalPatient{}).
		Where("phone_number = ? AND deleted_at IS NULL", phone).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("phone_number already exists")
	}

	rec := models.HospitalPatient{
		IDCardNumber: idcard,
		FirstName:    fn,
		LastName:     ln,
		PhoneNumber:  phone,
		BirthDay:     in.BirthDay, // เก็บเป็น DATE ตามโมเดล
		Gender:       g,
		CreatedAt:    now(),
		UpdatedAt:    now(),
	}
	if err := db.Create(&rec).Error; err != nil {
		return nil, err
	}

	res := dto.NewHospitalPatientResponse(rec)
	return &res, nil
}

// =====================
// LIST (no pagination)
// =====================

// คืนรายการทั้งหมด (หรือกรองด้วย q ถ้ามี)
// ค้นด้วย ILIKE ใน: id_card_number, phone_number, first_name, last_name
func ListHospitalPatients(db *gorm.DB, q string) ([]dto.HospitalPatientResponse, error) {
	var items []models.HospitalPatient

	tx := db.Model(&models.HospitalPatient{})
	if s := norm(q); s != "" {
		like := "%" + s + "%"
		tx = tx.Where(`
			id_card_number ILIKE ? OR
			phone_number  ILIKE ? OR
			first_name    ILIKE ? OR
			last_name     ILIKE ?`,
			like, like, like, like,
		)
	}

	// เรียงล่าสุดก่อน (id DESC) ตามสไตล์ที่ใช้บ่อย
	if err := tx.Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}

	out := make([]dto.HospitalPatientResponse, 0, len(items))
	for _, m := range items {
		out = append(out, dto.NewHospitalPatientResponse(m))
	}
	return out, nil
}

// =====================
// GET ONE
// =====================

func GetHospitalPatientByID(db *gorm.DB, id uint) (*dto.HospitalPatientResponse, error) {
	var m models.HospitalPatient
	if err := db.First(&m, id).Error; err != nil {
		return nil, err
	}
	res := dto.NewHospitalPatientResponse(m)
	return &res, nil
}

// =====================
// UPDATE
// =====================

func UpdateHospitalPatient(db *gorm.DB, id uint, in *dto.UpdateHospitalPatientDTO) (*dto.HospitalPatientResponse, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	var m models.HospitalPatient
	if err := db.First(&m, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]any{
		"updated_at": now(),
	}

	// id_card_number
	if in.IDCardNumber != nil {
		idc := onlyDigits(*in.IDCardNumber)
		if len(idc) != 13 {
			return nil, errors.New("id_card_number must be 13 digits")
		}
		// ตรวจซ้ำ ยกเว้นตัวเอง และต้องเป็นแถวที่ยังไม่ลบ
		var cnt int64
		if err := db.Model(&models.HospitalPatient{}).
			Where("id_card_number = ? AND id <> ? AND deleted_at IS NULL", idc, id).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errors.New("id_card_number already exists")
		}
		updates["id_card_number"] = idc
	}

	// phone_number
	if in.PhoneNumber != nil {
		ph := onlyDigits(*in.PhoneNumber)
		if len(ph) != 10 {
			return nil, errors.New("phone_number must be 10 digits")
		}
		var cnt int64
		if err := db.Model(&models.HospitalPatient{}).
			Where("phone_number = ? AND id <> ? AND deleted_at IS NULL", ph, id).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errors.New("phone_number already exists")
		}
		updates["phone_number"] = ph
	}

	// first_name / last_name
	if in.FirstName != nil {
		updates["first_name"] = norm(*in.FirstName)
	}
	if in.LastName != nil {
		updates["last_name"] = norm(*in.LastName)
	}

	// birth_day
	if in.BirthDay != nil {
		updates["birth_day"] = *in.BirthDay
	}

	// gender
	if in.Gender != nil {
		g := norm(*in.Gender)
		if !validGender(g) {
			return nil, errors.New(`gender must be "ชาย" or "หญิง"`)
		}
		updates["gender"] = g
	}

	// อัปเดต
	if err := db.Model(&m).Updates(updates).Error; err != nil {
		return nil, err
	}

	// โหลดล่าสุดคืน
	if err := db.First(&m, id).Error; err != nil {
		return nil, err
	}
	res := dto.NewHospitalPatientResponse(m)
	return &res, nil
}

// =====================
// DELETE (soft delete)
// =====================

func DeleteHospitalPatient(db *gorm.DB, id uint) error {
	return db.Delete(&models.HospitalPatient{}, id).Error
}
