package routes

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"strconv"
	"github.com/fouradithep/pillmate/models"
)

func SetupAdminAuthRoutes(app *fiber.App) {
	adm := app.Group("/admin") // อยู่นอก /api เพื่อไม่ชน Auth ของ mobile

	// POST /admin/login  -> set HttpOnly cookie: admin_jwt
	// {
	//   	"username": "admin_user",   // username ของแอดมิน/หมอในตาราง web_admins
	//   	"password": "secret123"     // เทียบกับ hash (bcrypt) ใน DB
	// }
	adm.Post("/login", func(c *fiber.Ctx) error {
		// ต้องเป็น JSON
		if ct := strings.ToLower(c.Get("Content-Type")); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "content-type must be application/json"})
		}
		// รับ input
		var in struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		in.Username = strings.TrimSpace(in.Username)
		in.Password = strings.TrimSpace(in.Password)
	if in.Username == "" || in.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username and password are required"})
		}

		// ตรวจ creds จาก DB เท่านั้น
		token, user, err := handlers.LoginAdmin(db.DB, in.Username, in.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}

		// ตั้งคุกกี้ JWT ฝั่งแอดมิน (ให้หมดอายุสอดคล้องกับ exp ใน JWT; ที่ฝั่ง handler ตั้งไว้ 24 ชม.)
		c.Cookie(&fiber.Cookie{
			Name:     "admin_jwt",
			Value:    token,
			Path:     "/",
			HTTPOnly: true,                       // JS อ่านไม่ได้ (กัน XSS)
			Secure:   true,                       // โปรดักชันบน HTTPS ให้เปิด true
			SameSite: "Lax",                      // หรือ "Strict" ถ้าใช้โดเมนเดียว
			Expires:  time.Now().Add(72 * time.Hour), // ให้ตรงกับ exp ใน JWT
			// Domain: "admin.yourdomain.com",     // ถ้ามีซับโดเมนให้ตั้ง
		})

		// ตอบกลับให้ front รู้ role สำหรับ redirect; ไม่จำเป็นต้องส่ง token ใน body แล้ว
		return c.JSON(fiber.Map{
			"message": "login success",
			"user": user,
			"role": user.Role, // เช่น "superadmin" | "doctor" | "staff"
		})
	})

	// POST /admin/logout  -> เคลียร์คุกกี้
	adm.Post("/logout", func(c *fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{
			Name:     "admin_jwt",
			Value:    "",
			Path:     "/",
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
			Expires:  time.Unix(0, 0), // หมดอายุย้อนหลัง
			// Domain: "admin.yourdomain.com",
		})
		return c.JSON(fiber.Map{
		"message": "logout success",
	})
	})

	// GET /admin/me -> ตรวจคุกกี้/เฮดเดอร์ ด้วย pickToken + parseJWT แล้วคืน {user, role}
	adm.Get("/me", func(c *fiber.Ctx) error {
		tok := handlers.PickToken(c) // ถ้า helpers เป็น lowercase ให้แก้ชื่อเรียกให้ตรง (pickToken)
		if strings.TrimSpace(tok) == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		claims, err := handlers.ParseJWT(tok) // เช่น parseJWT (แก้ชื่อเรียกให้ตรงกับที่ export)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		// admin_id จาก claims อาจเป็น float64 หรือ string ในบางเคส
		var adminID uint
		if v, ok := claims["admin_id"].(float64); ok {
			adminID = uint(v)
		} else if s, ok := claims["admin_id"].(string); ok && s != "" {
			if id64, _ := strconv.ParseUint(s, 10, 64); id64 > 0 {
				adminID = uint(id64)
			}
		}
		if adminID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		var user models.WebAdmin
		if err := db.DB.First(&user, adminID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		user.Password = "" // กันหลุด hash

		return c.JSON(fiber.Map{
			"user": user,
			"role": user.Role,
		})
	})
}
