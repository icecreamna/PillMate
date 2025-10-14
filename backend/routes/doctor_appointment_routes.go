package routes

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/handlers"
	"gorm.io/gorm"
)

// การใช้งาน:
//
// doctor := app.Group("/doctor/",
//     handlers.AuthAny,
//     handlers.RequireRole("doctor", "admin-app"),
// )
// routes.SetupDoctorAppointmentRoutes(doctor)
//
// เส้นทางที่ได้:
//   POST   /doctor/appointments
//   GET    /doctor/appointments?q=&date_from=&date_to=   (เรียงใหม่สุดก่อน)
//   GET    /doctor/appointments/:id
//   PUT    /doctor/appointments/:id
//   DELETE /doctor/appointments/:id
//
// หมายเหตุ: ส่ง Header Content-Type: application/json เมื่อมี Body

// --- helper: แปลง "YYYY-MM-DD" -> 00:00:00 @ UTC (ไม่อิงโซนไทย) ---
func parseQueryDateYMD_AsDateUTC(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("empty")
	}
	d, err := time.Parse("2006-01-02", s) // ไม่มี Location => UTC
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC), nil
}

func SetupDoctorAppointmentRoutes(api fiber.Router) {

	// CREATE
	api.Post("/appointments", func(c *fiber.Ctx) error {
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}
		var in dto.DoctorCreateAppointmentDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		docID, _ := c.Locals("admin_id").(uint)
		created, err := handlers.DoctorCreateAppointment(db.DB, &in, docID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "created", "data": created})
	})

	// LIST (เรียงใหม่สุดก่อน) + filter q/date_from/date_to
	// ตัวอย่าง: GET /doctor/appointments?q=11017&date_from=2025-10-01&date_to=2025-10-31
	api.Get("/appointments", func(c *fiber.Ctx) error {
		docID, _ := c.Locals("admin_id").(uint)
		q := c.Query("q", "")

		// parse date_from/date_to (YYYY-MM-DD -> UTC midnight)
		var dfPtr, dtPtr *time.Time
		if s := strings.TrimSpace(c.Query("date_from", "")); s != "" {
			if tUTC, err := parseQueryDateYMD_AsDateUTC(s); err == nil {
				dfPtr = &tUTC
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date_from (want YYYY-MM-DD)"})
			}
		}
		if s := strings.TrimSpace(c.Query("date_to", "")); s != "" {
			if tUTC, err := parseQueryDateYMD_AsDateUTC(s); err == nil {
				dtPtr = &tUTC
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date_to (want YYYY-MM-DD)"})
			}
		}

		list, err := handlers.DoctorListAppointments(db.DB, docID, q, dfPtr, dtPtr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": list})
	})

	// GET ONE (ของหมอคนนี้)
	api.Get("/appointments/:id", func(c *fiber.Ctx) error {
		docID, _ := c.Locals("admin_id").(uint)
		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		rec, err := handlers.DoctorGetAppointmentByID(db.DB, docID, uint(idU64))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": rec})
	})

	// UPDATE
	api.Put("/appointments/:id", func(c *fiber.Ctx) error {
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}
		docID, _ := c.Locals("admin_id").(uint)
		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		var in dto.DoctorUpdateAppointmentDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		updated, err := handlers.DoctorUpdateAppointment(db.DB, docID, uint(idU64), &in)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "updated", "data": updated})
	})

	// DELETE (soft) — หมอยกเลิกนัดของตนเอง
	api.Delete("/appointments/:id", func(c *fiber.Ctx) error {
		docID, _ := c.Locals("admin_id").(uint)
		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		if err := handlers.DoctorDeleteAppointment(db.DB, docID, uint(idU64)); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	})
}
