package routes

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/handlers"
)

func SetupNotifyRoutes(api fiber.Router) {

	// GET /api/notify/due-now?window=1
	api.Get("/notify/due-now", func(c *fiber.Ctx) error {
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"Unauthorized"})
		}
		w, _ := strconv.Atoi(c.Query("window", "1"))
		items, err := handlers.GetDueNow(db.DB, patientID, w)
		if err != nil { return c.Status(500).JSON(fiber.Map{"error": err.Error()}) }
		return c.JSON(fiber.Map{"data": dto.NotiItemsToDTO(items)})
	})

}
