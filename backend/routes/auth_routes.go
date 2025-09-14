package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
	"github.com/fouradithep/pillmate/db"
	"time"
)

func SetupAuthRoutes(app *fiber.App) {

	// (body: { "email": "test2@example.com", "password": "1234" })
	app.Post("/login", func(c *fiber.Ctx) error {
		patient := new(models.Patient)

		if err := c.BodyParser(patient); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		token, err := handlers.LoginPatient(db.DB, patient)

		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteNoneMode,  //fiber.CookieSameSiteNoneMode,
		})

		return c.JSON(fiber.Map{
			"message": "Login Successful",
			"token":   token, 
		})

	})

	app.Post("/logout", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",  // ชื่อ cookie ที่เคยเก็บ JWT token
			Value:    "",	// กำหนดให้เป็นค่าว่าง = ลบค่าเดิมออก
			Expires:  time.Now().Add(-time.Hour), // ตั้งเวลาให้หมดอายุไปแล้ว = ลบ cookie นี้
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteNoneMode,  //fiber.CookieSameSiteNoneMode,
		})
		return c.JSON(fiber.Map{"message": "Logged out"})
	})
}