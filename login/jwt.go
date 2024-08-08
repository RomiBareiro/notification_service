package login

import (
	"encoding/json"
	"net/http"
	t "notification_service/types"
	"time"

	"github.com/golang-jwt/jwt"
)

func (l *Login) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != "testuser" || password != "testpassword" { //TODO: improve user data adquisition
		errorResponse := t.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			l.Logger.Sugar().Errorf("could not encoder response: ", err)
		}
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(l.JWTKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(map[string]string{"token": tokenString}); err != nil {
		errorResponse := t.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Error encoding response",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err = json.NewEncoder(w).Encode(errorResponse); err != nil {
			l.Logger.Sugar().Errorf("could not encoder response: ", errorResponse)
		}
	}
}
