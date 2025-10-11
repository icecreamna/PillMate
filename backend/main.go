package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/routes"
)

func main() {
	// init DB + env (ตามที่ db.Init() จัดการโหลด .env แล้ว)
	db.Init()
	fmt.Println("Server started...")

	app := fiber.New()

	// ใช้โหมดควบคุมการเปิดเส้นทาง: "mobile" | "admin" | "all"
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "all" // ค่าเดิมให้เปิดทั้งสองฝั่ง เพื่อไม่กระทบพฤติกรรมเดิม
	}

	// อ่านพอร์ต (ดึงมาไว้ตรงนี้เพื่อพิมพ์ log ได้เลย)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// พิมพ์โหมดและพอร์ตให้เห็นชัดใน log
	fmt.Printf("APP_MODE=%s PORT=%s\n", mode, port)

	// ---------- Public (ไม่ต้องล็อกอิน) ----------
	// ใช้ app ตรง ๆ เลย (ฝั่ง Mobile)
	if mode == "mobile" || mode == "all" {
		routes.SetupPatientRoutes(app)
		routes.SetupOTPRoutes(app)
		routes.SetupForgotPasswordRoutes(app)
		routes.SetupPasswordRoutes(app)
		routes.SetupAuthRoutes(app)
		routes.SetupInitialDataRoutes(app)
	}

	// ---------- Protected (ต้องล็อกอิน) ----------
	// ทุกอย่างใต้ /api จะต้องผ่าน AuthRequired (ของ Mobile)
	if mode == "mobile" || mode == "all" {
		api := app.Group("/api", handlers.AuthRequired)

		routes.SetupMyMedicineRoutes(api)
		routes.SetupMedicineInfoRoutes(api)
		routes.SetupGroupMedicineRoutes(api)
		routes.SetupNotiInfosRoutes(api)
		routes.SetupNotiItemsRoutes(api)
		routes.SetupNotifyRoutes(api)
		routes.SetupProfileRoutes(api)
		routes.SetupSymptomRoutes(api)
		routes.SetupAppointmentRoutes(api)
	}

	// ---------- Admin/Web (เพิ่มโดยไม่กระทบ Mobile) ----------
	if mode == "admin" || mode == "all" {
		// สำคัญ: เส้นทางล็อกอินของเว็บให้อยู่ใต้ /admin (ไม่ใช่ /api)
		// - POST /admin/login  (Unified: superadmin หรือ doctor)
		routes.SetupAdminAuthRoutes(app)

		// กลุ่มสำหรับ superadmin จัดการ doctor (ต้องมี token + role superadmin)
		admin := app.Group("/admin/",
			handlers.AuthAny,
			handlers.RequireRole("superadmin", "admin-app"),
		)
		routes.SetupDoctorRoutes(admin)

		// ถ้ายังไม่มีหน้า panel ฝั่ง doctor ให้ไม่ต้องประกาศกลุ่มนี้ เพื่อลด unused var
		// doc := app.Group("/api/doctor", handlers.AuthAny, handlers.RequireRole("doctor", "admin-app"))
		// routes.SetupDoctorPanelRoutes(doc)
	}

	// start server
	if err := app.Listen(":" + port); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
