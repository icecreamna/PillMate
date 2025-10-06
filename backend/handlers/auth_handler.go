package handlers

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthRequired(c *fiber.Ctx) error {
	var tokenString string
	// ลองอ่านจาก cookie ก่อน 
	tokenString = c.Cookies("jwt")
	
	// ถ้าไม่มี cookie ให้อ่านจาก Authorization header 
	if tokenString == "" {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}
	// ถ้าไม่มี token เลย
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing token",
		})
	}

	jwtSecretKey := os.Getenv("jwtSecretKey") //in .env

	if jwtSecretKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "JWT secret key not set",
		})
	}

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
    return []byte(jwtSecretKey), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
    	"error": "Invalid token",
	})
	}

	claim := token.Claims.(jwt.MapClaims)

	// แปลง user_id เป็น uint หรือ int64 ก่อนเก็บ (โดยทั่วไป token จะเก็บเป็น float64)
	patientIDFloat, ok := claim["patient_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid patient_id in token"})
	}

	patientID := uint(patientIDFloat)
	c.Locals("patient_id", patientID)

	log.Printf("Authenticated patien ID: %d", patientID)


	return c.Next()
}

