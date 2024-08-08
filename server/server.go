package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"notification_service/login"
	"notification_service/service"
	"notification_service/types"
	t "notification_service/types"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server interface {
	SendHandler(w http.ResponseWriter, r *http.Request)
}

type server struct {
	Logger *zap.Logger
	Svc    *service.NotificationService
	ctx    context.Context
}

func NewServer(ctx context.Context, svc *service.NotificationService) *server {
	return &server{
		Svc:    svc,
		ctx:    ctx,
		Logger: svc.Logger,
	}
}
func ValidateInputData(body io.Reader) (types.InputInfo, error) {
	var in types.InputInfo
	err := json.NewDecoder(body).Decode(&in)
	if err != nil {
		return types.InputInfo{}, errors.New("invalid JSON format")
	}

	if in.Recipient == "" || in.NotificationGroup == "" {
		return types.InputInfo{}, errors.New("missing required fields")
	}

	in.NotificationGroup = types.NotificationType(strings.ToUpper(string(in.NotificationGroup)))
	if !types.IsValidNotificationType(in.NotificationGroup) {
		return types.InputInfo{}, fmt.Errorf("invalid notification type: %s", in.NotificationGroup)
	}

	return in, nil
}

func (s *server) SendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	in, err := ValidateInputData(r.Body)
	if err != nil {
		errorResponse := t.ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			s.Logger.Sugar().Errorf("could not encoder response: ", err)
		}
		return
	}

	s.Svc.Input = in
	out, err := s.Svc.SendNotification(s.ctx)
	if err != nil {
		errorResponse := t.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			s.Logger.Sugar().Errorf("could not encoder response: ", err)
		}
		return
	}

	response, err := json.Marshal(out)
	if err != nil {
		errorResponse := t.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			s.Logger.Sugar().Errorf("could not encoder response: ", err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		s.Logger.Error("could not send status OK")
	}
}

func ServerSetup(svc *service.NotificationService) *server {
	s := NewServer(context.Background(), svc)
	l := &login.Login{
		JWTKey: []byte("your_secret_key"),
		Ctx:    s.ctx,
		Logger: s.Logger,
	}
	router := mux.NewRouter()

	router.HandleFunc("/login", l.LoginHandler).Methods("POST")

	protectedRoutes := router.PathPrefix("/V1").Subrouter()
	protectedRoutes.Use(l.ValidateJWTMiddleware)
	protectedRoutes.HandleFunc("/notify", s.SendHandler).Methods("POST")

	port := ":8080"
	s.Logger.Sugar().Infof("Listening port: %s", port)
	s.Logger.Sugar().Fatal(http.ListenAndServe(port, router))

	return s
}
