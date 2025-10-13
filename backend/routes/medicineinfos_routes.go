package routes

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

// หมายเหตุ: ให้เรียกฟังก์ชันนี้ภายใต้กลุ่มที่มี middleware แล้ว เช่น
// admin := app.Group("/admin/", handlers.AuthAny, handlers.RequireRole("superadmin", "admin-app"))
// routes.SetupMedicineInfoRoutes(admin)
//
// เส้นทางจริงจะเป็น:
//   POST   /admin/medicine-info
//   GET    /admin/medicine-infos
//   GET    /admin/medicine-info/:id
//   PUT    /admin/medicine-info/:id
//   DELETE /admin/medicine-info/:id
func SetupMedicineInfoRoutes(api fiber.Router) {
	// CREATE
	// Body ตัวอย่าง:
	// {
	//   "med_name": "Paracetamol",
	//   "generic_name": "Acetaminophen",
	//   "properties": "แก้ปวด ลดไข้",
	//   "strength": "500 mg",
	//   "form_id": 1,
	//   "unit_id": 1,
	//   "instruction_id": 2,
	//   "med_status": "active"
	// }
	api.Post("/medicine-info", func(c *fiber.Ctx) error {
		// บังคับเป็น JSON
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).
				JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}

		var in models.MedicineInfo
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		created, err := handlers.AddMedicineInfo(db.DB, &in)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "created",
			"data":    dto.ToMedicineInfoDTO(created),
		})
	})

	// LIST
	api.Get("/medicine-infos", func(c *fiber.Ctx) error {
		list, err := handlers.GetMedicineInfos(db.DB)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": dto.ToMedicineInfoDTOs(list)})
	})

	// GET ONE
	api.Get("/medicine-info/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		m, err := handlers.GetMedicineInfo(db.DB, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": dto.ToMedicineInfoDTO(m)})
	})

	// UPDATE (อัปเดตเฉพาะฟิลด์ที่ส่งมา)
	api.Put("/medicine-info/:id", func(c *fiber.Ctx) error {
		// บังคับเป็น JSON
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).
				JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}

		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		var in models.MedicineInfo
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		updated, err := handlers.UpdateMedicineInfo(db.DB, uint(id), &in)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "updated", "data": dto.ToMedicineInfoDTO(updated)})
	})

	// DELETE
	api.Delete("/medicine-info/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		if err := handlers.DeleteMedicineInfo(db.DB, uint(id)); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	})
}

// =======================
// Read-only for doctor
// ใช้งานภายใต้กลุ่มที่มี middleware แล้ว เช่น
// doctor := app.Group("/doctor/",
//     handlers.AuthAny,
//     handlers.RequireRole("doctor", "admin-app"),
// )
// routes.SetupDoctorMedicineReadRoutes(doctor)
//
// เส้นทางจริงจะเป็น:
//   GET /doctor/medicine-infos
//   GET /doctor/medicine-info/:id
// =======================
func SetupDoctorMedicineReadRoutes(api fiber.Router) {
	// LIST
	api.Get("/medicine-infos", func(c *fiber.Ctx) error {
		list, err := handlers.GetMedicineInfos(db.DB)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": dto.ToMedicineInfoDTOs(list)})
	})

	// GET ONE
	api.Get("/medicine-info/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		m, err := handlers.GetMedicineInfo(db.DB, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": dto.ToMedicineInfoDTO(m)})
	})
}