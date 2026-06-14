package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/example/go-user-api/internal/logger"
)

// RequestLogger logs method, path, status, and latency for every request.
func RequestLogger() fiber.Handler {
	log := logger.Get()
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		log.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", c.IP()),
		)
		return err
	}
}

// RequestID injects a unique X-Request-ID header into every response.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Use Fiber's built-in locals mechanism; generate a simple ID if absent.
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = generateID()
		}
		c.Set("X-Request-ID", reqID)
		c.Locals("requestID", reqID)
		return c.Next()
	}
}

// generateID produces a lightweight pseudo-unique ID without importing uuid.
func generateID() string {
	return strconvItoa(time.Now().UnixNano())
}

func strconvItoa(n int64) string {
	// tiny helper – avoids importing strconv just for this
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 20)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
