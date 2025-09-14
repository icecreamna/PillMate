package routes

import (
	"errors"
	"strconv"
	"github.com/fouradithep/pillmate/models"

	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"gorm.io/gorm"
)

func SetupMedicineInfoRoutes(api fiber.Router) {
	// CREATE
	// Body
	// {
	// "med_name": "Paracetamol",
	// "generic_name": "Acetaminophen",
	// "properties": "แก้ปวด ลดไข้",
	// "strength": "500 mg",
	// "form_id": 1,
	// "unit_id": 1,
	// "instruction_id": 2,
	// "med_status": "active"
	// }
	// POST /api/medicine-info
	api.Post("/medicine-info", func(c *fiber.Ctx) error {
		var in models.MedicineInfo
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		created, err := handlers.AddMedicineInfo(db.DB, &in)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "created", "data": created})
	})

	// LIST
	// GET /api/medicine-infos
	api.Get("/medicine-infos", func(c *fiber.Ctx) error {
		list, err := handlers.GetMedicineInfos(db.DB)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": list})
	})

	// GET ONE
	// GET /api/medicine-info/:id
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
		return c.JSON(fiber.Map{"data": m})
	})

	// UPDATE
	// PUT /api/medicine-info/:id
	api.Put("/medicine-info/:id", func(c *fiber.Ctx) error {
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
		return c.JSON(fiber.Map{"message": "updated", "data": updated})
	})

	// DELETE
	// DELETE /api/medicine-info/:id
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