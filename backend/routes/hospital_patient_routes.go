package routes

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/handlers"
	"gorm.io/gorm"
)

// การใช้งาน (แนะนำให้วางไว้ใน main/routes setup):
//
// // กลุ่มสำหรับ doctor (ต้องมี token + role doctor หรือ admin-app)
// doctor := app.Group("/doctor/",
//     handlers.AuthAny,
//     handlers.RequireRole("doctor", "admin-app"),
// )
// routes.SetupHospitalPatientRoutes(doctor)
//
// เมื่อใช้งานแบบนี้ เส้นทางจะเป็น:
//   POST   /doctor/hospital-patients
//   GET    /doctor/hospital-patients?q=...
//   GET    /doctor/hospital-patients/:id
//   PUT    /doctor/hospital-patients/:id
//   DELETE /doctor/hospital-patients/:id

func SetupHospitalPatientRoutes(api fiber.Router) {
	// =========================
	// CREATE
	// POST /doctor/hospital-patients
	// {
	//   "id_card_number": "1234567890123",
	//   "first_name": "Somchai",
	//   "last_name": "Jaidee",
	//   "phone_number": "0812345678",
	//   "birth_day": "2000-01-15T00:00:00+07:00",
	//   "gender": "ชาย"
	// }
	// =========================
	api.Post("/hospital-patients", func(c *fiber.Ctx) error {
		// บังคับ Content-Type เป็น JSON
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).
				JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}

		var in dto.CreateHospitalPatientDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		created, err := handlers.CreateHospitalPatient(db.DB, &in)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "created",
			"data":    created,
		})
	})

	// =========================
	// LIST (ไม่แบ่งหน้า) + ค้นหา q
	// GET /doctor/hospital-patients?q=
	// =========================
	api.Get("/hospital-patients", func(c *fiber.Ctx) error {
		q := c.Query("q", "")

		items, err := handlers.ListHospitalPatients(db.DB, q)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": items})
	})

	// =========================
	// GET ONE
	// GET /doctor/hospital-patients/:id
	// =========================
	api.Get("/hospital-patients/:id", func(c *fiber.Ctx) error {
		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		rec, err := handlers.GetHospitalPatientByID(db.DB, uint(idU64))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": rec})
	})

	// =========================
	// UPDATE (partial)
	// PUT /doctor/hospital-patients/:id
	// Body (UpdateHospitalPatientDTO): ฟิลด์ที่ไม่ส่ง = ไม่แก้
	// =========================
	api.Put("/hospital-patients/:id", func(c *fiber.Ctx) error {
		// บังคับ Content-Type เป็น JSON
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).
				JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}

		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		var in dto.UpdateHospitalPatientDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		updated, err := handlers.UpdateHospitalPatient(db.DB, uint(idU64), &in)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"message": "updated",
			"data":    updated,
		})
	})

	// =========================
	// DELETE (soft delete ตามโมเดล)
	// DELETE /doctor/hospital-patients/:id
	// =========================
	api.Delete("/hospital-patients/:id", func(c *fiber.Ctx) error {
		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		if err := handlers.DeleteHospitalPatient(db.DB, uint(idU64)); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	})
}
