package middleware

import (
	"fmt"
	"strings"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

)

// Logger middleware logs HTTP request details
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Store the time in locals
		c.Locals("start_time", start)

		// Continue with the next middleware/handler
		err := c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Log the details
		fmt.Printf(
			"[%s] %s %s %s - %v\n",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Method(),
			c.Path(),
			c.IP(),
			duration,
		)

		return err
	}
}

// FUNC AUTHMIDDLEWARE CHECKS FOR VALID JWT TOKEN
func AuthMiddleware(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorzation")

		// Check if auth header exists and has bearer token
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Missing Authorzation header",
			})
		}

		// Split bearer token
		parts := strings.Split(authHeader, "")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected singing method: %v", token.Header["alg"])

			}

			return []byte(secretKey), nil
		})

		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Check if token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Store user information in context
			c.Locals("user_id", claims["user_id"])
			c.Locals("role", claims["role"])
			return c.Next()
		}

		return c.Status(401).JSON(fiber.Map{
			"error": "invalid token claims",
		})
	}
}

// ROLEAUTH MIDDLEWARE CHECKS IF USER HAS REQUIRED ROLE
func RoleAuth(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != requiredRole {
			return c.Status(403).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// REQUESTVALIDATOR VALIDATES REQUIRED FIELDS IN THE REQUEST BODY
func RequestValidator(requiredFields []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := make(map[string]interface{})

		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": " Invalid request body",
			})

		}

		// Check for required fields
		for _, field := range requiredFields {
			if _, exists := body[field]; !exists {
				return c.Status(400).JSON(fiber.Map{
					"error": fmt.Sprintf("Missing required field: %s", field),
				})
			}
		}

		return c.Next()
	}
}

// RATELIMITER IMPLEMENT A SMILE RATE LIMITING MIDDLEWARE
func RateLimiter(request int, duration time.Duration) fiber.Handler {
	// Store for rate limiting
	type client struct {
		count    int
		lastSeen time.Time
	}
	clients := make(map[string]*client)

	return func(c *fiber.Ctx) error {
		ip := c.IP()
		now := time.Now()

		// Get or create client record
		cl, exists := clients[ip]
		if !exists {
			clients[ip] = &client{count: 0, lastSeen: now}
			cl = clients[ip]
		}

		// Reset count if duration has passed
		if now.Sub(cl.lastSeen) > duration {
			cl.count = 0
			cl.lastSeen = now
		}

		if cl.count >= request {
			return c.Status(429).JSON(fiber.Map{
				"error": "Rate limit exceeded", 
			})
		}

		// Update count and timestamp
		cl.count++
		cl.lastSeen = now
		return c.Next()
	}

}
