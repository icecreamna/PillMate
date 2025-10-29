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
// doctor := app.Group("/doctor",
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
//
// หมายเหตุ response จะมี patient_code เสมอ (ระบบ gen อัตโนมัติ)

func SetupHospitalPatientRoutes(api fiber.Router) {
	isJSON := func(ct string) bool {
		ct = strings.ToLower(ct)
		return strings.Contains(ct, "application/json") || strings.Contains(ct, "+json")
	}
	asErr := func(code int, msg string) error {
		return fiber.NewError(code, msg)
	}

	// =========================
	// CREATE
	// POST /doctor/hospital-patients
	// Request (CreateHospitalPatientDTO):
	// {
	//   "id_card_number": "1234567890123",
	//   "first_name": "Somchai",
	//   "last_name": "Jaidee",
	//   "phone_number": "0812345678",
	//   "birth_day": "2000-01-15"  // หรือ "2000-01-15T00:00:00+07:00"
	//   "gender": "ชาย"
	// }
	// Response: 201
	// {
	//   "message": "created",
	//   "data": {
	//     "id": 1,
	//     "patient_code": "000001",
	//     "id_card_number": "1234567890123",
	//     "first_name": "Somchai",
	//     "last_name": "Jaidee",
	//     "phone_number": "0812345678",
	//     "birth_day": "2000-01-15",
	//     "gender": "ชาย",
	//     "created_at": "...",
	//     "updated_at": "..."
	//   }
	// }
	// =========================
	api.Post("/hospital-patients", func(c *fiber.Ctx) error {
		if ct := c.Get("Content-Type"); !isJSON(ct) {
			return asErr(fiber.StatusUnsupportedMediaType, "Content-Type must be application/json or *+json")
		}

		var in dto.CreateHospitalPatientDTO
		if err := c.BodyParser(&in); err != nil {
			return asErr(fiber.StatusBadRequest, "invalid request body")
		}

		created, err := handlers.CreateHospitalPatient(db.DB, &in)
		if err != nil {
			msg := err.Error()
			switch {
			case strings.Contains(msg, "missing required fields"),
				strings.Contains(msg, "must be 13 digits"),
				strings.Contains(msg, "must be 10 digits"),
				strings.Contains(msg, `gender must be`):
				return asErr(fiber.StatusUnprocessableEntity, msg)
			case strings.Contains(msg, "already exists"):
				return asErr(fiber.StatusConflict, msg)
			default:
				return asErr(fiber.StatusBadRequest, msg)
			}
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "created",
			"data":    created,
		})
	})

	// =========================
	// LIST (ไม่แบ่งหน้า) + ค้นหา q
	// GET /doctor/hospital-patients?q=
	// หมายเหตุ: ฝั่ง handler ค้นหาได้ทั้ง patient_code, id_card_number, phone_number, first_name, last_name
	// =========================
	api.Get("/hospital-patients", func(c *fiber.Ctx) error {
		q := c.Query("q", "")

		items, err := handlers.ListHospitalPatients(db.DB, q)
		if err != nil {
			return asErr(fiber.StatusInternalServerError, err.Error())
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
			return asErr(fiber.StatusBadRequest, "invalid id")
		}

		rec, err := handlers.GetHospitalPatientByID(db.DB, uint(idU64))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return asErr(fiber.StatusNotFound, "not found")
			}
			return asErr(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(fiber.Map{"data": rec})
	})

	// =========================
	// UPDATE (partial)
	// PUT /doctor/hospital-patients/:id
	// Body (UpdateHospitalPatientDTO): ฟิลด์ที่ไม่ส่ง = ไม่แก้
	// =========================
	api.Put("/hospital-patients/:id", func(c *fiber.Ctx) error {
		if ct := c.Get("Content-Type"); !isJSON(ct) {
			return asErr(fiber.StatusUnsupportedMediaType, "Content-Type must be application/json or *+json")
		}

		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return asErr(fiber.StatusBadRequest, "invalid id")
		}

		var in dto.UpdateHospitalPatientDTO
		if err := c.BodyParser(&in); err != nil {
			return asErr(fiber.StatusBadRequest, "invalid request body")
		}

		updated, err := handlers.UpdateHospitalPatient(db.DB, uint(idU64), &in)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return asErr(fiber.StatusNotFound, "not found")
			}
			msg := err.Error()
			switch {
			case strings.Contains(msg, "must be 13 digits"),
				strings.Contains(msg, "must be 10 digits"),
				strings.Contains(msg, `gender must be`):
				return asErr(fiber.StatusUnprocessableEntity, msg)
			case strings.Contains(msg, "already exists"):
				return asErr(fiber.StatusConflict, msg)
			default:
				return asErr(fiber.StatusBadRequest, msg)
			}
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
			return asErr(fiber.StatusBadRequest, "invalid id")
		}

		if err := handlers.DeleteHospitalPatient(db.DB, uint(idU64)); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return asErr(fiber.StatusNotFound, "not found")
			}
			return asErr(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	})
}
