package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// RequestID middleware config
type RequestIDConfig struct {
	// Next defines a function to skip this middleware when returned true.
	Next func(c fiber.Ctx) bool

	// Header is the header key where to get/set the unique request ID
	// Default: "X-Request-ID"
	Header string

	// Generator defines a function to generate the unique identifier.
	// Default: uuid.NewString()
	Generator func() string
}

// ConfigDefault is the default config
var ConfigDefault = RequestIDConfig{
	Next:      nil,
	Header:    fiber.HeaderXRequestID,
	Generator: uuid.NewString,
}

// Helper function to set default values
func configDefault(config ...RequestIDConfig) RequestIDConfig {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Header == "" {
		cfg.Header = ConfigDefault.Header
	}

	if cfg.Generator == nil {
		cfg.Generator = ConfigDefault.Generator
	}

	return cfg
}

// NewRequestID creates a new middleware handler
func NewRequestID(config ...RequestIDConfig) fiber.Handler {
	cfg := configDefault(config...)

	return func(c fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Get id from request, else we generate one
		rid := c.Get(cfg.Header)
		if rid == "" {
			rid = cfg.Generator()
		}

		// Set new id to response header
		c.Set(cfg.Header, rid)

		// Add the request ID to locals
		c.Locals("requestid", rid)

		// Continue stack
		return c.Next()
	}
}
