package web

import (
	"errors"
	"net/http"

	appErrors "github.com/LalatinaHub/LatinaServer/pkg/errors"
	"github.com/LalatinaHub/LatinaServer/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents standard error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Type    string `json:"type,omitempty"`
}

// ErrorMiddleware handles centralized error responses in Gin
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors accumulated
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var appErr *appErrors.AppError
			if errors.As(err, &appErr) {
				logger.Error().
					Str("type", string(appErr.Type)).
					Str("message", appErr.Message).
					Err(appErr.Err).
					Msg("API error handled")

				status := http.StatusInternalServerError
				switch appErr.Type {
				case appErrors.ErrorTypeValidation:
					status = http.StatusBadRequest
				case appErrors.ErrorTypeDatabase:
					status = http.StatusServiceUnavailable
				case appErrors.ErrorTypeNetwork:
					status = http.StatusBadGateway
				case appErrors.ErrorTypeConfig:
					status = http.StatusInternalServerError
				}

				c.JSON(status, ErrorResponse{
					Success: false,
					Error:   appErr.Message,
					Type:    string(appErr.Type),
				})
				return
			}

			// Generic error
			logger.Error().Err(err).Msg("Unhandled API error")
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Error:   "Internal server error",
				Type:    string(appErrors.ErrorTypeInternal),
			})
		}
	}
}
