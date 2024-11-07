package service

import (
	"context"
	"fmt"
	"time"

	d "notification_service/db"
	t "notification_service/types"

	"go.uber.org/zap"
)

type Service interface {
	SendNotification(ctx context.Context, input t.InputInfo) (t.Output, error)
	IsAllowed(ctx context.Context, nType t.NotificationType) bool
}

type NotificationService struct {
	DB                  *d.DBConnector
	Logger              *zap.Logger
	CurrentNotification t.Notifications
	Input               t.InputInfo
}

func NewNotificationService(logger *zap.Logger, conn *d.DBConnector) *NotificationService {
	return &NotificationService{Logger: logger, DB: conn}
}

func (s *NotificationService) SendNotification(ctx context.Context) (t.Output, error) {
	if !s.IsAllowed(ctx, s.Input.NotificationGroup) {
		return t.Output{}, fmt.Errorf("rate limit exceeded for recipient %s", s.Input.Recipient)
	}

	if err := s.DB.RecordNotification(ctx, &s.Input, s.CurrentNotification); err != nil {
		return t.Output{}, fmt.Errorf("could not update notifications table: %v", err)
	}
	if err := sendTelegramMessage(5751493884, "HOLA ROMI ENVIADO!"); err != nil {
		return t.Output{}, fmt.Errorf("could not send telegram message: %v", err)
	}
	// TODO: Implement sending logic here
	s.Logger.Info("notification is sent",
		zap.String("recipient", s.Input.Recipient),
		zap.String("type", string(s.Input.NotificationGroup)),
	)
	output := t.Output{
		Recipient:         s.Input.Recipient,
		NotificationGroup: s.Input.NotificationGroup,
		TimeStamp:         time.Now(),
		Message:           "Notification sent successfully",
	}

	return output, nil
}

func (s *NotificationService) IsAllowed(ctx context.Context, nType t.NotificationType) bool {
	lastNotification := make(chan t.Notifications)
	rateLimitRules := make(chan t.RateLimitRule)

	errChan := make(chan error)

	go func() {
		n, err := s.DB.GetLastNotification(ctx, s.Input)
		if err != nil {
			errChan <- err
			return
		}
		lastNotification <- n

		r, err := s.DB.GetRateLimitRules(ctx, nType)
		if err != nil {
			errChan <- err
			return
		}
		rateLimitRules <- r
	}()
	// Handle go routines results
	var n t.Notifications
	var r t.RateLimitRule
	var err error

	for i := 0; i < 2; i++ {
		select {
		case n = <-lastNotification:
		case r = <-rateLimitRules:
		case err = <-errChan:
			s.Logger.Error("Error retrieving data", zap.Error(err))
			return false
		}
	}

	if n == (t.Notifications{}) || n.ID == "" {
		return true // first notif
	}

	if isGreaterDuration(r.Duration, n.CreatedAt) {
		return true // notif is allowed to be sent because time to send new notif of this type is reached
	}

	if r.MaxCount > n.Counter { // if current notif < maxcount send message & incr counter
		s.CurrentNotification = n
		s.CurrentNotification.Counter = n.Counter + 1
		return true
	}

	return false
}

// isGreaterDuration returns true if timeSince creation is greater than the rate limit duration
func isGreaterDuration(durationSeconds float64, creation time.Time) bool {
	now := time.Now().UTC()
	timeSince := now.Sub(creation).Seconds()
	v := timeSince > durationSeconds
	return v
}
