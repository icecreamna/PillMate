package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
)

// ใช้คู่กับ handlers.GetPatient และ handlers.UpdatePatientBasic
func SetupProfileRoutes(api fiber.Router) {

	// READ — โปรไฟล์ของตัวเอง
	// GET /api/patient/me
	api.Get("/patient/me", func(c *fiber.Ctx) error {
		// auth
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		patient, err := handlers.GetPatient(db.DB, patientID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": patient})
	})

	// UPDATE — อัปเดตเฉพาะ 4 ฟิลด์ที่อนุญาต
	// PUT /api/patient/me
	// Body (อัปเดตเฉพาะฟิลด์ที่ส่งมา เป็นค่าว่างไม่อัปเดต):
	// {
	//   "id_card_number": "1234567890123",
	//   "first_name": "Somchai",
	//   "last_name": "Jaidee",
	//   "phone_number": "0812345678"
	// }
	api.Put("/patient/me", func(c *fiber.Ctx) error {
		// auth
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var payload models.Patient
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		updated, err := handlers.UpdatePatientBasic(db.DB, patientID, &payload)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"message": "updated",
			"data":    updated,
		})
	})
}
