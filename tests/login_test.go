package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	l "notification_service/login"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestValidateJWTMiddleware(t *testing.T) {
	// Create a valid token
	claims := &jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
	}
	// Create the login instance with the middleware
	l := &l.Login{
		JWTKey: []byte("your_secret_key"),
		Ctx:    context.Background(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(l.JWTKey)

	// Create a handler that the middleware will wrap
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name               string
		authHeader         string
		expectedStatusCode int
	}{
		{
			name:               "Valid Token",
			authHeader:         "Bearer " + tokenString,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid Token",
			authHeader:         "Bearer invalid_token",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Missing Token",
			authHeader:         "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/V1/notify", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			middleware := l.ValidateJWTMiddleware(handler)
			middleware.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
