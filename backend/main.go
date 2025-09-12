package main

import (
	"github.com/fouradithep/pillmate/db"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/routes"
	"os"
)

func main() {
	db.Init()
	fmt.Println("Server started...")
	

	app := fiber.New()

	routes.SetupPatientRoutes(app)
	routes.SetupOTPRoutes(app)
	routes.SetupForgotPasswordRoutes(app)
	routes.SetupPasswordRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}
	if err := app.Listen(":" + port); err != nil {
		fmt.Println("Error starting server:", err)
	}

}