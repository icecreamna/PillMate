package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/routes"
)

func main() {
	// init DB + env (db.Init() จะโหลด .env ให้)
	db.Init()
	fmt.Println("Server started...")

	app := fiber.New()

	// CORS: ใส่ origin ของหน้าเว็บที่เรียกจริง (ตัวอย่าง Vite dev)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		ExposeHeaders:    "Content-Type",
		AllowCredentials: true, // สำคัญ! เมื่อใช้คุกกี้
	}))

	// อ่านพอร์ต (พิมพ์ log ให้เห็น)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("PORT=%s\n", port)

	// ---------- Public (ไม่ต้องล็อกอิน) ----------
	// ฝั่ง Mobile (public endpoints)
	routes.SetupPatientRoutes(app)
	routes.SetupOTPRoutes(app)
	routes.SetupForgotPasswordRoutes(app)
	routes.SetupPasswordRoutes(app)
	routes.SetupAuthRoutes(app)
	routes.SetupInitialDataRoutes(app)

	// ---------- Protected (ต้องล็อกอิน) ----------
	// กลุ่ม /api ของ Mobile ต้องผ่าน AuthRequired
	api := app.Group("/api", handlers.AuthRequired)
	routes.SetupMyMedicineRoutes(api)
	routes.SetupGroupMedicineRoutes(api)
	routes.SetupNotiInfosRoutes(api)
	routes.SetupNotiItemsRoutes(api)
	routes.SetupNotifyRoutes(api)
	routes.SetupProfileRoutes(api)
	routes.SetupSymptomRoutes(api)
	routes.SetupMobileAppointmentRoutes(api)
	

	// ---------- Admin/Web ----------------------------
	routes.SetupAdminAuthRoutes(app)
	routes.SetupDoctorSelfRoutes(app)

	// กลุ่มสำหรับ superadmin 
	admin := app.Group("/admin",
		handlers.AuthAny,
		handlers.RequireRole("superadmin"),
	)
	routes.SetupDoctorRoutes(admin)
	routes.SetupMedicineInfoRoutes(admin)

	// กลุ่มสำหรับ doctor
	doctor := app.Group("/doctor",
		handlers.AuthAny,
		handlers.RequireRole("doctor"),
	)
	routes.SetupHospitalPatientRoutes(doctor)
	routes.SetupPrescriptionRoutes(doctor)
	routes.SetupDoctorAppointmentRoutes(doctor)
	routes.SetupDoctorMedicineReadRoutes(doctor)
	routes.SetupDoctorPublicRoutes(doctor)


	// start server
	if err := app.Listen(":" + port); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
