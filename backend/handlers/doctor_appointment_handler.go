package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/models"
)

const defaultHospitalID uint = 1 // โรงพยาบาลเดียว (seed ไว้ id=1)

// ===== Fixed timezone (Bangkok) สำหรับการตีความ "เวลา" (ไม่ใช่ "วัน") =====
var handlerBangkokLoc = func() *time.Location {
	l, _ := time.LoadLocation("Asia/Bangkok")
	return l
}()

// ============== DATE PARSER (ไม่อิงโซนไทย) =======================
// parse "YYYY-MM-DD" -> ตรึงเป็นเที่ยงคืน UTC ของวันนั้น (เพื่อให้ DB/GET ตรงวันเดิม)
func parseDateYMD_AsDateUTC(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("appointment_date is required")
	}
	d, err := time.Parse("2006-01-02", s) // parse แบบไม่ใส่ Location = time.UTC
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: want YYYY-MM-DD, got %q", s)
	}
	// เก็บเป็น 00:00:00 UTC ของวันนั้น
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC), nil
}

// ============== TIME PARSER (อิงโซนไทย) =========================
// parse "HH:MM" (or "HH:MM:SS") -> anchor 2000-01-01 HH:MM @Bangkok → UTC เพื่อเซฟ
func parseTimeHMToUTC(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("appointment_time is required")
	}
	var t time.Time
	var err error
	// รองรับ HH:MM และ HH:MM:SS
	t, err = time.ParseInLocation("15:04", s, handlerBangkokLoc)
	if err != nil {
		t, err = time.ParseInLocation("15:04:05", s, handlerBangkokLoc)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid time format: want HH:MM, got %q", s)
		}
	}
	// ใช้ปีสมัยใหม่ (2000) เพื่อกัน offset ประวัติศาสตร์ +06:42/BC
	anchor := time.Date(2000, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, handlerBangkokLoc)
	return anchor.UTC(), nil
}

// ============ CREATE (หมอสร้างนัด) ============
func DoctorCreateAppointment(db *gorm.DB, in *dto.DoctorCreateAppointmentDTO, doctorIDFromToken uint) (*dto.DoctorAppointmentResponse, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	idc := OnlyDigits(in.IDCardNumber)
	if len(idc) != 13 {
		return nil, errors.New("id_card_number must be 13 digits")
	}

	// parse date ("YYYY-MM-DD") เป็น "UTC midnight ของวันนั้น" (ไม่อิงโซนไทย)
	dateUTC, err := parseDateYMD_AsDateUTC(in.AppointmentDate)
	if err != nil {
		return nil, err
	}
	// parse time ("HH:MM") เป็น anchor 2000-01-01 @Bangkok แล้วแปลงเป็น UTC
	timeUTC, err := parseTimeHMToUTC(in.AppointmentTime)
	if err != nil {
		return nil, err
	}

	// resolve hospital_id (optional -> default=1)
	hospID := defaultHospitalID
	if in.HospitalID != nil && *in.HospitalID != 0 {
		hospID = *in.HospitalID
	}

	// resolve doctor_id (optional -> token)
	docID := doctorIDFromToken
	if in.DoctorID != nil && *in.DoctorID != 0 {
		docID = *in.DoctorID
	}
	if docID == 0 {
		return nil, errors.New("doctor_id is required")
	}

	// FK checks
	{
		var h models.Hospital
		if err := db.First(&h, hospID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("hospital not found")
			}
			return nil, err
		}
	}
	{
		var doc models.WebAdmin
		if err := db.Where("role = ?", "doctor").First(&doc, docID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("doctor not found")
			}
			return nil, err
		}
	}

	// กันชน slot หมอ (doctor_id + date + time) — ใช้ค่าที่จะเซฟจริง (UTC)
	{
		var cnt int64
		if err := db.Model(&models.Appointment{}).
			Where("doctor_id = ? AND appointment_date = ? AND appointment_time = ?", docID, dateUTC, timeUTC).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errors.New("this doctor already has an appointment at the given date/time")
		}
	}

	rec := models.Appointment{
		IDCardNumber:    idc,
		AppointmentDate: dateUTC, // UTC midnight ของวันนั้น (เพื่อให้ DB โชว์วันตรงกับอินพุต)
		AppointmentTime: timeUTC, // UTC (anchor 2000-01-01)
		HospitalID:      hospID,
		DoctorID:        docID,
		Note:            in.Note,
	}
	if err := db.Create(&rec).Error; err != nil {
		return nil, err
	}

	res := dto.ToDoctorAppointmentResponse(rec)
	return &res, nil
}

// ============ LIST (ของหมอคนนั้น) — เรียงใหม่สุดก่อน ============
func DoctorListAppointments(db *gorm.DB, doctorID uint, q string, dateFrom, dateTo *time.Time) ([]dto.DoctorAppointmentResponse, error) {
	var rows []models.Appointment
	tx := db.Model(&models.Appointment{}).Where("doctor_id = ?", doctorID)

	if s := Norm(q); s != "" {
		like := "%" + s + "%"
		tx = tx.Where("id_card_number ILIKE ?", like)
	}
	// หมายเหตุ: แนะนำให้ตัวเรียก (routes) ส่งค่า df/dt เป็น "UTC midnight" ของวันนั้นแล้ว
	if dateFrom != nil && !dateFrom.IsZero() {
		tx = tx.Where("appointment_date >= ?", *dateFrom)
	}
	if dateTo != nil && !dateTo.IsZero() {
		tx = tx.Where("appointment_date <= ?", *dateTo)
	}

	if err := tx.Order("appointment_date DESC, appointment_time DESC, id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return dto.ToDoctorAppointmentResponses(rows), nil
}

// ============ GET ONE (เช็คว่าเป็นของหมอคนนี้) ============
func DoctorGetAppointmentByID(db *gorm.DB, doctorID, id uint) (*dto.DoctorAppointmentResponse, error) {
	var rec models.Appointment
	if err := db.First(&rec, id).Error; err != nil {
		return nil, err
	}
	if rec.DoctorID != doctorID {
		return nil, gorm.ErrRecordNotFound // กันรั่วข้อมูลคนอื่น
	}
	res := dto.ToDoctorAppointmentResponse(rec)
	return &res, nil
}

// ============ UPDATE (partial) — อนุญาตหมอแก้ของตนเอง ============
func DoctorUpdateAppointment(db *gorm.DB, doctorID, id uint, in *dto.DoctorUpdateAppointmentDTO) (*dto.DoctorAppointmentResponse, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	var rec models.Appointment
	if err := db.First(&rec, id).Error; err != nil {
		return nil, err
	}
	if rec.DoctorID != doctorID {
		return nil, gorm.ErrRecordNotFound
	}

	updates := map[string]any{
		"updated_at": time.Now().UTC(), // ใช้ UTC ให้คงที่
	}

	// id_card_number
	if in.IDCardNumber != nil {
		idc := OnlyDigits(*in.IDCardNumber)
		if len(idc) != 13 {
			return nil, errors.New("id_card_number must be 13 digits")
		}
		updates["id_card_number"] = idc
	}

	// เตรียมค่าที่จะใช้กันชน slot หลังอัปเดต (เริ่มจากค่าปัจจุบัน)
	newDate := rec.AppointmentDate
	newTime := rec.AppointmentTime

	// appointment_date: string "YYYY-MM-DD" → ตีความเป็น "UTC midnight ของวันนั้น"
	if in.AppointmentDate != nil {
		dUTC, err := parseDateYMD_AsDateUTC(*in.AppointmentDate)
		if err != nil {
			return nil, err
		}
		updates["appointment_date"] = dUTC
		newDate = dUTC
	}

	// appointment_time: string "HH:MM" → UTC (anchor 2000-01-01)
	if in.AppointmentTime != nil {
		tUTC, err := parseTimeHMToUTC(*in.AppointmentTime)
		if err != nil {
			return nil, err
		}
		updates["appointment_time"] = tUTC
		newTime = tUTC
	}

	// hospital (optional override)
	if in.HospitalID != nil {
		if *in.HospitalID == 0 {
			return nil, errors.New("hospital_id is required")
		}
		var h models.Hospital
		if err := db.First(&h, *in.HospitalID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("hospital not found")
			}
			return nil, err
		}
		updates["hospital_id"] = *in.HospitalID
	}

	// doctor (เปลี่ยนเจ้าของนัดได้หรือไม่? ปกติไม่ — ถ้าจะอนุญาตให้ตรวจ role และความถูกต้อง)
	if in.DoctorID != nil && *in.DoctorID != 0 && *in.DoctorID != rec.DoctorID {
		var doc models.WebAdmin
		if err := db.Where("role = ?", "doctor").First(&doc, *in.DoctorID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("doctor not found")
			}
			return nil, err
		}
		updates["doctor_id"] = *in.DoctorID
	}

	// note (nullable)
	if in.Note != nil {
		updates["note"] = in.Note
	}

	// กันชน slot ใหม่ (ใช้ doctor_id หลังอัปเดตถ้ามี)
	newDoctor := rec.DoctorID
	if v, ok := updates["doctor_id"].(uint); ok && v != 0 {
		newDoctor = v
	}
	{
		var cnt int64
		if err := db.Model(&models.Appointment{}).
			Where("doctor_id = ? AND appointment_date = ? AND appointment_time = ? AND id <> ?",
				newDoctor, newDate, newTime, rec.ID).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errors.New("this doctor already has an appointment at the given date/time")
		}
	}

	if err := db.Model(&rec).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := db.First(&rec, id).Error; err != nil {
		return nil, err
	}
	res := dto.ToDoctorAppointmentResponse(rec)
	return &res, nil
}

// ============ DELETE (soft) — หมอยกเลิกนัดของตนเอง ============
func DoctorDeleteAppointment(db *gorm.DB, doctorID, id uint) error {
	var rec models.Appointment
	if err := db.Select("id, doctor_id").First(&rec, id).Error; err != nil {
		return err
	}
	if rec.DoctorID != doctorID {
		return gorm.ErrRecordNotFound
	}
	return db.Delete(&models.Appointment{}, id).Error
}
