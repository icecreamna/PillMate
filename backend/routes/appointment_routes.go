package routes

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
)

func SetupAppointmentRoutes(api fiber.Router) {

	// =========================================================
	// GET ONE
	// GET /api/appointment/:id
	// อิงสิทธิ์จาก patient_id -> โหลด id_card_number แล้วคิวรีด้วย (id, id_card_number)
	// =========================================================
	api.Get("/appointment/:id", func(c *fiber.Ctx) error {
		// ⛑️ auth
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// แปลงพารามิเตอร์ id
		appointmentIDU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || appointmentIDU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		appointmentID := uint(appointmentIDU64)

		// โหลดเลขบัตรประชาชนจากผู้ใช้ที่ล็อกอิน
		var me models.Patient
		if err := db.DB.Select("id_card_number").Where("id = ?", patientID).First(&me).Error; err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrRecordNotFound) {
				status = fiber.StatusUnauthorized
			}
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		if me.IDCardNumber == nil || strings.TrimSpace(*me.IDCardNumber) == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing id_card_number for this account",
			})
		}
		idCard := strings.TrimSpace(*me.IDCardNumber)

		// ดึงใบนัดหนึ่งรายการ (ยืนยันด้วย id + id_card_number)
		appointment, getErr := handlers.GetAppointment(db.DB, appointmentID, idCard)
		if getErr != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(getErr, gorm.ErrRecordNotFound) {
				status = fiber.StatusNotFound
			}
			return c.Status(status).JSON(fiber.Map{"error": getErr.Error()})
		}
		return c.JSON(fiber.Map{"data": dto.AppointmentToDTO(*appointment)})
	})

	// =========================================================
	// LIST
	// GET /api/appointments?date_from=YYYY-MM-DD&date_to=YYYY-MM-DD&hospital_id=&doctor_id=
	// หมายเหตุ: ใช้ id_card_number จากผู้ใช้ที่ล็อกอินเพื่อจำกัดเรคคอร์ด
	// =========================================================
	api.Get("/appointments", func(c *fiber.Ctx) error {
		// ⛑️ auth
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// โหลดเลขบัตรประชาชนจากผู้ใช้ที่ล็อกอิน
		var me models.Patient
		if err := db.DB.Select("id_card_number").Where("id = ?", patientID).First(&me).Error; err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrRecordNotFound) {
				status = fiber.StatusUnauthorized
			}
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		if me.IDCardNumber == nil || strings.TrimSpace(*me.IDCardNumber) == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing id_card_number for this account",
			})
		}

		idCard := strings.TrimSpace(*me.IDCardNumber)

		// อ่าน query params
		dateFromStr := strings.TrimSpace(c.Query("date_from"))
		dateToStr := strings.TrimSpace(c.Query("date_to"))
		hospitalIDStr := strings.TrimSpace(c.Query("hospital_id"))
		doctorIDStr := strings.TrimSpace(c.Query("doctor_id"))

		const ymd = "2006-01-02"

		q := db.DB.Model(&models.Appointment{}).
			Preload("Hospital").
			Preload("WebAdmin").
			Where("id_card_number = ?", idCard)

		// กรองช่วงวันที่ (ถ้าส่งมา)
		if dateFromStr != "" {
			if t, err := time.ParseInLocation(ymd, dateFromStr, time.Local); err == nil {
				q = q.Where("appointment_date >= ?", t)
			}
		}
		if dateToStr != "" {
			if t, err := time.ParseInLocation(ymd, dateToStr, time.Local); err == nil {
				q = q.Where("appointment_date <= ?", t)
			}
		}

		// กรองด้วย hospital_id / doctor_id (ถ้าส่งมา)
		if hospitalIDStr != "" {
			if hid, err := strconv.ParseUint(hospitalIDStr, 10, 64); err == nil {
				q = q.Where("hospital_id = ?", uint(hid))
			}
		}
		if doctorIDStr != "" {
			if did, err := strconv.ParseUint(doctorIDStr, 10, 64); err == nil {
				q = q.Where("doctor_id = ?", uint(did))
			}
		}

		var appointments []models.Appointment
		if err := q.Order("appointment_date, appointment_time").Find(&appointments).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"data": dto.AppointmentsToDTO(appointments)})
	})
}
