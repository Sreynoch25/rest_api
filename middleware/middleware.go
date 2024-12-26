package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Protected() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get token from header
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Missing authorization header",
            })
        }
        
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        // Parse token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte("secret"), nil
        })

        // Check if there is an error
        if err != nil || !token.Valid {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid or expired token",
            })
        }

        // Add claims to context
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid token claims",
            })
        }

        c.Locals("user", claims)
        return c.Next()
    }
}