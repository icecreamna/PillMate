package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"
	// "unicode"

	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/models"
)

// ===== ค่าพื้นฐาน =====
const defaultHospitalID uint = 1 // โรงพยาบาลเดียว (seed ไว้ id=1)

// ===== Fixed timezone (Bangkok) สำหรับตีความ "เวลา HH:MM" =====
var handlerBangkokLoc = func() *time.Location {
	l, _ := time.LoadLocation("Asia/Bangkok")
	return l
}()

// ===== Helpers =====

// // เก็บเฉพาะตัวเลข (normalize id_card_number)
// func OnlyDigits(s string) string {
// 	var b strings.Builder
// 	for _, r := range s {
// 		if unicode.IsDigit(r) {
// 			b.WriteRune(r)
// 		}
// 	}
// 	return b.String()
// }

// // ตัดช่องว่าง + lower (ไว้ทำ ILIKE %...%)
// func Norm(s string) string {
// 	return strings.ToLower(strings.TrimSpace(s))
// }

// ===== DATE/TIME PARSERS =====

// "YYYY-MM-DD" -> 00:00:00 UTC (date-only ไม่เพี้ยน timezone)
func parseDateYMD_AsDateUTC(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("appointment_date is required")
	}
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: want YYYY-MM-DD, got %q", s)
	}
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC), nil
}

// "HH:MM" (หรือ "HH:MM:SS") -> anchor 2000-01-01 @Bangkok แล้ว .UTC() เพื่อเซฟ
func parseTimeHMToUTC(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("appointment_time is required")
	}
	tm, err := time.ParseInLocation("15:04", s, handlerBangkokLoc)
	if err != nil {
		tm, err = time.ParseInLocation("15:04:05", s, handlerBangkokLoc)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid time format: want HH:MM, got %q", s)
		}
	}
	anchor := time.Date(2000, 1, 1, tm.Hour(), tm.Minute(), tm.Second(), 0, handlerBangkokLoc)
	return anchor.UTC(), nil
}

// ============================================================================
// CREATE — ใช้ id_card_number จาก DTO เพื่อหา Patient -> เซ็ต patient_id + id_card_number แล้วบันทึก
// ============================================================================
func DoctorCreateAppointment(
	db *gorm.DB,
	in *dto.DoctorCreateAppointmentDTO,
	doctorIDFromToken uint,
) (*dto.DoctorAppointmentResponse, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	// หา Patient ด้วยเลขบัตรจาก DTO
	idc := OnlyDigits(in.IDCardNumber)
	if len(idc) != 13 {
		return nil, errors.New("id_card_number must be 13 digits")
	}
	var p models.Patient
	if err := db.Select("id").
		Where("id_card_number = ?", idc).
		First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, err
	}

	// parse วัน/เวลา
	dateUTC, err := parseDateYMD_AsDateUTC(in.AppointmentDate)
	if err != nil {
		return nil, err
	}
	timeUTC, err := parseTimeHMToUTC(in.AppointmentTime)
	if err != nil {
		return nil, err
	}

	// resolve hospital/doctor
	hospID := defaultHospitalID
	if in.HospitalID != nil && *in.HospitalID != 0 {
		hospID = *in.HospitalID
	}
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

	// กันซ้อนฝั่งหมอ (doctor + date + time)
	{
		var cnt int64
		if err := db.Model(&models.Appointment{}).
			Where("doctor_id = ? AND appointment_date = ? AND appointment_time = ?",
				docID, dateUTC, timeUTC).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errors.New("this doctor already has an appointment at the given date/time")
		}
	}

	// กันซ้อนฝั่งคนไข้ (patient + date + time)
	{
		var cnt int64
		if err := db.Model(&models.Appointment{}).
			Where("patient_id = ? AND appointment_date = ? AND appointment_time = ?",
				p.ID, dateUTC, timeUTC).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errors.New("this patient already has an appointment at the given date/time")
		}
	}

	// บันทึก (ใช้ idc ที่ normalize แล้ว ไม่อ้าง pointer)
	rec := models.Appointment{
		PatientID:       p.ID,
		IDCardNumber:    idc,
		AppointmentDate: dateUTC,
		AppointmentTime: timeUTC,
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

// ============================================================================
// LIST — ของหมอคนนั้น (order ใหม่สุดก่อน) + filter q/date_from/date_to
// q ค้นด้วย id_card_number ในตาราง appointments ได้เลย
// ============================================================================
func DoctorListAppointments(
	db *gorm.DB,
	doctorID uint,
	q string,
	dateFrom, dateTo *time.Time,
) ([]dto.DoctorAppointmentResponse, error) {
	var rows []models.Appointment

	tx := db.Model(&models.Appointment{}).
		Where("doctor_id = ?", doctorID)

	if s := Norm(q); s != "" {
		like := "%" + s + "%"
		tx = tx.Where("id_card_number ILIKE ?", like)
	}

	if dateFrom != nil && !dateFrom.IsZero() {
		tx = tx.Where("appointment_date >= ?", *dateFrom)
	}
	if dateTo != nil && !dateTo.IsZero() {
		tx = tx.Where("appointment_date <= ?", *dateTo)
	}

	if err := tx.Order("appointment_date DESC, appointment_time DESC, id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return dto.ToDoctorAppointmentResponses(rows), nil
}

// ============================================================================
// GET ONE — ของหมอคนนั้นเท่านั้น
// ============================================================================
func DoctorGetAppointmentByID(
	db *gorm.DB,
	doctorID, id uint,
) (*dto.DoctorAppointmentResponse, error) {
	var rec models.Appointment
	if err := db.Where("id = ? AND doctor_id = ?", id, doctorID).
		First(&rec).Error; err != nil {
		return nil, err
	}
	res := dto.ToDoctorAppointmentResponse(rec)
	return &res, nil
}

// ============================================================================
// UPDATE (partial) — ของหมอคนนั้นเท่านั้น
// - ถ้าส่ง id_card_number มา: หา Patient แล้วอัปเดต patient_id + id_card_number ตามนั้น
// - ไม่ส่ง: ไม่เปลี่ยนเจ้าของนัด
// ============================================================================
func DoctorUpdateAppointment(
	db *gorm.DB,
	doctorID, id uint,
	in *dto.DoctorUpdateAppointmentDTO,
) (*dto.DoctorAppointmentResponse, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	var rec models.Appointment
	if err := db.Where("id = ? AND doctor_id = ?", id, doctorID).
		First(&rec).Error; err != nil {
		return nil, err
	}

	updates := map[string]any{
		"updated_at": time.Now().UTC(),
	}

	// เปลี่ยนเจ้าของนัดด้วยเลขบัตร (ถ้าส่งมา)
	newPatient := rec.PatientID
	if in.IDCardNumber != nil {
		s := OnlyDigits(strings.TrimSpace(*in.IDCardNumber))
		if len(s) != 13 {
			return nil, errors.New("id_card_number must be 13 digits")
		}
		var p models.Patient
		if err := db.Select("id").
			Where("id_card_number = ?", s).
			First(&p).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("patient not found")
			}
			return nil, err
		}
		updates["patient_id"] = p.ID
		updates["id_card_number"] = s // ใช้ค่าที่ normalize แล้ว
		newPatient = p.ID
	}

	// เตรียมค่าที่ใช้กันชน slot หลังอัปเดต
	newDate := rec.AppointmentDate
	newTime := rec.AppointmentTime
	newDoctor := rec.DoctorID

	// appointment_date
	if in.AppointmentDate != nil {
		dUTC, err := parseDateYMD_AsDateUTC(*in.AppointmentDate)
		if err != nil {
			return nil, err
		}
		updates["appointment_date"] = dUTC
		newDate = dUTC
	}

	// appointment_time
	if in.AppointmentTime != nil {
		tUTC, err := parseTimeHMToUTC(*in.AppointmentTime)
		if err != nil {
			return nil, err
		}
		updates["appointment_time"] = tUTC
		newTime = tUTC
	}

	// hospital (optional)
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

	// doctor (ถ้าจะเปลี่ยนเจ้าของนัด)
	if in.DoctorID != nil && *in.DoctorID != 0 && *in.DoctorID != rec.DoctorID {
		var doc models.WebAdmin
		if err := db.Where("role = ?", "doctor").First(&doc, *in.DoctorID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("doctor not found")
			}
			return nil, err
		}
		updates["doctor_id"] = *in.DoctorID
		newDoctor = *in.DoctorID
	}

	// note (nullable)
	if in.Note != nil {
		updates["note"] = in.Note
	}

	// กันชน slot: ฝั่งหมอ
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

	// กันชน slot: ฝั่งคนไข้
	{
		var cnt int64
		if err := db.Model(&models.Appointment{}).
			Where("patient_id = ? AND appointment_date = ? AND appointment_time = ? AND id <> ?",
				newPatient, newDate, newTime, rec.ID).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errors.New("this patient already has an appointment at the given date/time")
		}
	}

	// บันทึก
	if err := db.Model(&rec).Updates(updates).Error; err != nil {
		return nil, err
	}
	// reload
	if err := db.First(&rec, rec.ID).Error; err != nil {
		return nil, err
	}

	res := dto.ToDoctorAppointmentResponse(rec)
	return &res, nil
}

// ============================================================================
// DELETE — ของหมอคนนั้นเท่านั้น (soft delete)
// ============================================================================
func DoctorDeleteAppointment(
	db *gorm.DB,
	doctorID, id uint,
) error {
	var rec models.Appointment
	if err := db.Select("id").
		Where("id = ? AND doctor_id = ?", id, doctorID).
		First(&rec).Error; err != nil {
		return err
	}
	return db.Delete(&models.Appointment{}, rec.ID).Error
}
