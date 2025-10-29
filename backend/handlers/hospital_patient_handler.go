package handlers

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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

// เริ่มเลข patient code จาก ENV (PATIENT_CODE_START) ถ้าไม่ตั้ง ใช้ 1
func patientCodeStart() int {
	v := strings.TrimSpace(os.Getenv("PATIENT_CODE_START"))
	if v == "" {
		return 1
	}
	if n, err := strconv.Atoi(v); err == nil && n > 0 {
		return n
	}
	return 1
}

// ต้องถูกเรียก "ภายใน Transaction" เท่านั้น
// ล็อคด้วย advisory lock เพื่อลด race จากหลายอินสแตนซ์
func nextPatientCode(tx *gorm.DB) (string, error) {
	// เลขคีย์ล็อกเป็นอะไรก็ได้ที่สม่ำเสมอในโปรเซสนี้
	if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", 424242).Error; err != nil {
		return "", err
	}

	// ดึงค่าสูงสุดของ patient_code (ความยาวคงที่ 6 ตัวอักษร เปรียบได้)
	var maxCode string
	if err := tx.
		Raw(`SELECT COALESCE(MAX(patient_code), '') FROM hospital_patients WHERE deleted_at IS NULL`).
		Scan(&maxCode).Error; err != nil {
		return "", err
	}

	start := patientCodeStart()
	var next int
	if maxCode == "" {
		next = start
	} else {
		d := maxCode
		if len(d) > 6 {
			d = d[len(d)-6:]
		}
		n, err := strconv.Atoi(d)
		if err != nil {
			// ถ้า parse ไม่ได้ fallback เป็น start
			next = start
		} else {
			next = n + 1
			if next < start {
				next = start
			}
		}
	}

	return fmt.Sprintf("%06d", next), nil
}

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

	var out dto.HospitalPatientResponse

	// ใช้ทรานแซกชันเพื่อให้ออก patient_code แบบ atomic + กัน race
	err := db.Transaction(func(tx *gorm.DB) error {
		// ตรวจซ้ำ (respect soft delete)
		var count int64
		if err := tx.Model(&models.HospitalPatient{}).
			Where("id_card_number = ? AND deleted_at IS NULL", idcard).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("id_card_number already exists")
		}
		if err := tx.Model(&models.HospitalPatient{}).
			Where("phone_number = ? AND deleted_at IS NULL", phone).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("phone_number already exists")
		}

		// ออกเลข patient_code
		code, err := nextPatientCode(tx)
		if err != nil {
			return err
		}

		rec := models.HospitalPatient{
			PatientCode:  code,
			IDCardNumber: idcard,
			FirstName:    fn,
			LastName:     ln,
			PhoneNumber:  phone,
			BirthDay:     in.BirthDay, // เก็บเป็น DATE ตามโมเดล
			Gender:       g,
			CreatedAt:    now(),
			UpdatedAt:    now(),
		}

		// ลอง insert; ถ้าชน unique patient_code (เคสหายาก) จะ retry ออกรหัสใหม่
		const maxTry = 2
		for try := 0; try <= maxTry; try++ {
			if err := tx.Create(&rec).Error; err != nil {
				// ตรวจข้อความ error ชน index ของ patient_code (ชื่อ index ต้องตรงกับโมเดล)
				if strings.Contains(err.Error(), "uniq_patient_code_active") && try < maxTry {
					if code, err = nextPatientCode(tx); err != nil {
						return err
					}
					rec.PatientCode = code
					continue
				}
				return err
			}
			break // สำเร็จ
		}

		r := dto.NewHospitalPatientResponse(rec)
		out = r
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &out, nil
}

// =====================
// LIST (no pagination)
// =====================

// คืนรายการทั้งหมด (หรือกรองด้วย q ถ้ามี)
// ค้นด้วย ILIKE ใน: patient_code, id_card_number, phone_number, first_name, last_name
func ListHospitalPatients(db *gorm.DB, q string) ([]dto.HospitalPatientResponse, error) {
	var items []models.HospitalPatient

	tx := db.Model(&models.HospitalPatient{})
	if s := norm(q); s != "" {
		like := "%" + s + "%"
		tx = tx.Where(`
			patient_code  ILIKE ? OR
			id_card_number ILIKE ? OR
			phone_number  ILIKE ? OR
			first_name    ILIKE ? OR
			last_name     ILIKE ?`,
			like, like, like, like, like,
		)
	}

	// เรียงล่าสุดก่อน (id DESC)
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
