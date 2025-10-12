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

// parse "YYYY-MM-DD" -> time.Time (date-only)
func parseDateYMD(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("appointment_date is required")
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: want YYYY-MM-DD, got %q", s)
	}
	return t, nil
}

// parse "HH:MM" (or "HH:MM:SS") -> time.Time (time-only)
func parseTimeHM(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("appointment_time is required")
	}
	if t, err := time.Parse("15:04", s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("15:04:05", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid time format: want HH:MM, got %q", s)
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

	// parse date ("YYYY-MM-DD") and time ("HH:MM")
	dateOnly, err := parseDateYMD(in.AppointmentDate)
	if err != nil {
		return nil, err
	}
	tOnly, err := parseTimeHM(in.AppointmentTime)
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

	// กันชน slot หมอ (doctor_id + date + time)
	{
		var cnt int64
		if err := db.Model(&models.Appointment{}).
			Where("doctor_id = ? AND appointment_date = ? AND appointment_time = ?", docID, dateOnly, tOnly).
			Count(&cnt).Error; err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, errors.New("this doctor already has an appointment at the given date/time")
		}
	}

	rec := models.Appointment{
		IDCardNumber:    idc,
		AppointmentDate: dateOnly,
		AppointmentTime: tOnly,
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
		"updated_at": time.Now(),
	}

	// id_card_number
	if in.IDCardNumber != nil {
		idc := OnlyDigits(*in.IDCardNumber)
		if len(idc) != 13 {
			return nil, errors.New("id_card_number must be 13 digits")
		}
		updates["id_card_number"] = idc
	}

	// เตรียมค่าที่จะใช้กันชน slot หลังอัปเดต
	newDate := rec.AppointmentDate
	newTime := rec.AppointmentTime

	// appointment_date: string "YYYY-MM-DD"
	if in.AppointmentDate != nil {
		d, err := parseDateYMD(*in.AppointmentDate)
		if err != nil {
			return nil, err
		}
		updates["appointment_date"] = d
		newDate = d
	}

	// appointment_time: string "HH:MM"
	if in.AppointmentTime != nil {
		tOnly, err := parseTimeHM(*in.AppointmentTime)
		if err != nil {
			return nil, err
		}
		updates["appointment_time"] = tOnly
		newTime = tOnly
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
