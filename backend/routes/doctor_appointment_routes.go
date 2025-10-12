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

func SetupDoctorAppointmentRoutes(api fiber.Router) {

	// CREATE
	// ตัวอย่าง body:
	// {
	//   "id_card_number": "1101700203452",
	//   "appointment_date": "2025-10-20",
	//   "appointment_time": "14:30",
	//   "note": "งดอาหารก่อนตรวจ 8 ชั่วโมง"
	// }
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

		// parse date_from/date_to (YYYY-MM-DD)
		var dfPtr, dtPtr *time.Time
		if s := strings.TrimSpace(c.Query("date_from", "")); s != "" {
			if t, err := time.Parse("2006-01-02", s); err == nil {
				tt := t // local date
				dfPtr = &tt
			}
		}
		if s := strings.TrimSpace(c.Query("date_to", "")); s != "" {
			if t, err := time.Parse("2006-01-02", s); err == nil {
				tt := t
				dtPtr = &tt
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
	// ตัวอย่าง body:
	// {
	//   "appointment_date": "2025-10-25",
	//   "appointment_time": "10:00",
	//   "note": "เลื่อนเวลา"
	// }
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
