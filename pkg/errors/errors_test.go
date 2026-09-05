package errors_test

import (
	"errors"
	"testing"

	appErrors "github.com/LalatinaHub/LatinaServer/pkg/errors"
)

func TestAppError(t *testing.T) {
	baseErr := errors.New("underlying failure")

	tests := []struct {
		name       string
		err        *appErrors.AppError
		targetType appErrors.ErrorType
		hasUnder   bool
	}{
		{
			name:       "DatabaseError",
			err:        appErrors.NewDatabaseError("db query failed", baseErr),
			targetType: appErrors.ErrorTypeDatabase,
			hasUnder:   true,
		},
		{
			name:       "ConfigError",
			err:        appErrors.NewConfigError("config unmarshal failed", baseErr),
			targetType: appErrors.ErrorTypeConfig,
			hasUnder:   true,
		},
		{
			name:       "NetworkError",
			err:        appErrors.NewNetworkError("network request failed", baseErr),
			targetType: appErrors.ErrorTypeNetwork,
			hasUnder:   true,
		},
		{
			name:       "ValidationError",
			err:        appErrors.NewValidationError("invalid param", nil),
			targetType: appErrors.ErrorTypeValidation,
			hasUnder:   false,
		},
		{
			name:       "InternalError",
			err:        appErrors.NewInternalError("internal issue", baseErr),
			targetType: appErrors.ErrorTypeInternal,
			hasUnder:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Type != tt.targetType {
				t.Fatalf("expected type %v, got %v", tt.targetType, tt.err.Type)
			}

			if tt.hasUnder && !errors.Is(tt.err, baseErr) && errors.Unwrap(tt.err) != baseErr {
				t.Fatalf("expected unwrapped error to be %v", baseErr)
			}

			if tt.err.Error() == "" {
				t.Fatal("expected non-empty error string")
			}
		})
	}
}
