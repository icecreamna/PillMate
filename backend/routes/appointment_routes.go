package routes

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"gorm.io/gorm"
)

// จะได้เส้นทาง:
//   GET /api/appointments/latest   -> ใบนัดล่าสุดของผู้ใช้ที่ล็อกอิน
//   GET /api/appointments/:id      -> ใบนัดตาม id (ต้องเป็นของผู้ใช้เอง)
func SetupMobileAppointmentRoutes(api fiber.Router) {
	// GET /api/appointments/latest
	api.Get("/appointments/latest", func(c *fiber.Ctx) error {
		// ปิด cache responses (ข้อมูลส่วนบุคคล)
		c.Set("Cache-Control", "no-store")

		patientID, _ := c.Locals("patient_id").(uint) // ต้องมีจาก middleware auth ฝั่ง mobile
		if patientID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		latest, err := handlers.MobileGetLatestAppointment(db.DB, patientID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no appointment"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": latest})
	})

	// GET /api/appointments/:id
	api.Get("/appointments/:id", func(c *fiber.Ctx) error {
		// ปิด cache responses (ข้อมูลส่วนบุคคล)
		c.Set("Cache-Control", "no-store")

		patientID, _ := c.Locals("patient_id").(uint)
		if patientID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		rec, err := handlers.MobileGetAppointmentByID(db.DB, patientID, uint(idU64))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": rec})
	})
}
