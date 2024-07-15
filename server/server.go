package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"notification_service/service"
	"notification_service/types"

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

	if !types.IsValidNotificationType(in.NotificationGroup) {
		return types.InputInfo{}, fmt.Errorf("invalid notification type: %s", in.NotificationGroup)
	}
	return in, nil
}

func (s *server) SendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
		return
	}
	in, err := ValidateInputData(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	s.Svc.Input = in
	out, err := s.Svc.SendNotification(s.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	response, err := json.Marshal(out)
	if err != nil {
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
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
	http.HandleFunc("/sendNotif", s.SendHandler)

	// start server
	port := ":8080"
	s.Logger.Sugar().Infof("Listening port: %s", port)
	s.Logger.Sugar().Fatal(http.ListenAndServe(port, nil))

	return s

}