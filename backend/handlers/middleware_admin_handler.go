package handlers

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// ดึง token จาก cookie/header (รองรับทั้ง jwt และ admin_jwt)
func pickToken(c *fiber.Ctx) string {
	// if t := c.Cookies("jwt"); t != "" { // เผื่อบางที่ใช้ชื่อเดิม
	// 	return t
	// }
	if t := c.Cookies("admin_jwt"); t != "" {
		return t
	}
	ah := c.Get("Authorization")
	if ah != "" && strings.HasPrefix(ah, "Bearer ") {
		return strings.TrimPrefix(ah, "Bearer ")
	}
	return ""
}

func parseJWT(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("jwtSecretKey")
	if secret == "" {
		return nil, errors.New("JWT secret key not set")
	}
	parser := &jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Alg()}}
	tok, err := parser.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || tok == nil || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	claims, _ := tok.Claims.(jwt.MapClaims)

	// ตรวจ exp ถ้ามี
	if v, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(v) {
			return nil, errors.New("token expired")
		}
	}
	return claims, nil
}

// ใช้กับฝั่งเว็บ (admin/doctor) — ไม่กระทบ AuthRequired เดิมของ mobile
func AuthAny(c *fiber.Ctx) error {
	t := pickToken(c)
	if t == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}
	claims, err := parseJWT(t)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	// map claims -> locals (มีอะไรมาก็ค่อยๆ ใส่ลงไป)
	if v, ok := claims["role"].(string); ok { c.Locals("role", v) }
	if v, ok := claims["aud"].(string); ok { c.Locals("aud", v) }
	if v, ok := claims["admin_id"].(float64); ok { c.Locals("admin_id", uint(v)) }
	if v, ok := claims["doctor_id"].(float64); ok { c.Locals("doctor_id", uint(v)) }
	return c.Next()
}

// การ์ดสิทธิ์ตาม role (และออปชัน aud)
func RequireRole(required string, allowedAud ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != required {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		if len(allowedAud) > 0 {
			aud, _ := c.Locals("aud").(string)
			ok := false
			for _, a := range allowedAud {
				if a == aud { ok = true; break }
			}
			if !ok {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "invalid audience"})
			}
		}
		return c.Next()
	}
}
