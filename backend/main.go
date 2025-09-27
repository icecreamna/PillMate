package main

import (
	"github.com/fouradithep/pillmate/db"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/routes"
	"os"
	"github.com/fouradithep/pillmate/handlers"
)

func main() {
	db.Init()
	fmt.Println("Server started...")
	

	app := fiber.New()

	// ---------- Public (ไม่ต้องล็อกอิน) ----------
	// ใช้ app ตรง ๆ เลย
	routes.SetupPatientRoutes(app)
	routes.SetupOTPRoutes(app)
	routes.SetupForgotPasswordRoutes(app)
	routes.SetupPasswordRoutes(app)
	routes.SetupAuthRoutes(app)

	// ---------- Protected (ต้องล็อกอิน) ----------
	// ทุกอย่างใต้ /api จะต้องผ่าน AuthRequired
	api := app.Group("/api", handlers.AuthRequired)

	routes.SetupMyMedicineRoutes(api)
	routes.SetupMedicineInfoRoutes(api)
	routes.SetupGroupMedicineRoutes(api)
	

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}
	if err := app.Listen(":" + port); err != nil {
		fmt.Println("Error starting server:", err)
	}

}