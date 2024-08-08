package login

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	t "notification_service/types"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type Login struct {
	JWTKey []byte
	Ctx    context.Context
	Logger *zap.Logger
}
type contextKey string

const (
	claimsKey contextKey = "claims"
)

func (l *Login) ValidateJWTMiddleware(next http.Handler) http.Handler {
	var err error
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			err = fmt.Errorf("Authorization token is missing")
			errorResponse := t.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if err = json.NewEncoder(w).Encode(errorResponse); err != nil {
				l.Logger.Sugar().Errorf("could not encoder response: ", err)
			}
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, claims, err := validateToken(tokenString, l.JWTKey)
		if err != nil || !token.Valid {
			err = fmt.Errorf("Invalid token")
			errorResponse := t.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if err = json.NewEncoder(w).Encode(errorResponse); err != nil {
				l.Logger.Sugar().Errorf("could not encoder response: ", err)
			}
			return
		}

		// Add the claims to the request context
		ctx := context.WithValue(r.Context(), claimsKey, claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func validateToken(tokenString string, jwtKey []byte) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}
