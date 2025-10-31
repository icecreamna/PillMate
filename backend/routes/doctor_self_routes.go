package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/handlers"
)

// ====== Public self-register for Doctor ======
// ใช้ใน main.go โซน Public: routes.SetupDoctorSelfRoutes(app)
func SetupDoctorSelfRoutes(app fiber.Router) {
	// POST /auth/doctor/register
	// Body (CreateDoctorDTO):
	// {
	//   "username": "doc1",
	//   "password": "secret",
	//   "first_name": "Somchai",
	//   "last_name": "Jaidee"
	// }
	app.Post("/auth/doctor/register", func(c *fiber.Ctx) error {
		// บังคับให้เป็น JSON
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).
				JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}

		var in dto.CreateDoctorDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		// สมัครเอง actorID = 0  (handlers.CreateDoctor จะบังคับ role="doctor")
		created, err := handlers.CreateDoctor(db.DB, &in, 0)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "created",
			"data":    dto.ToWebAdminDTO(created), // ไม่คืน password
		})
	})
}
