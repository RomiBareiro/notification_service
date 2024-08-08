package types

import (
	"time"
)

type NotificationType string

const (
	Status    NotificationType = "STATUS"
	News      NotificationType = "NEWS"
	Marketing NotificationType = "MARKETING"
)

type RateLimitRule struct {
	ID               string           `db:"id"`
	NotificationType NotificationType `db:"notification_type"`
	MaxCount         int              `db:"max_count"`
	Duration         float64          `db:"duration"`
}
type Notifications struct {
	ID               string           `db:"id"`
	NotificationType NotificationType `db:"notification_type"`
	Recipient        string           `db:"recipient"`
	Counter          int              `db:"counter"`
	CreatedAt        time.Time        `db:"created_at"`
	UpdatedAt        time.Time        `db:"updated_at"`
}

type InputInfo struct {
	Recipient         string           `json:"recipient"  validate:"required"`
	NotificationGroup NotificationType `json:"group"  validate:"required"`
}
type Output struct {
	Recipient         string           `json:"recipient"`
	NotificationGroup NotificationType `json:"group"`
	TimeStamp         time.Time        `json:"timestamp,omitempty"`
	Message           string           `json:"message,omitempty"`
}

// harcoded by the moment, could be given by argument flags or secrets etc
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var ValidNotificationTypes = []NotificationType{Status, Marketing, News}

func IsValidNotificationType(t NotificationType) bool {
	for _, validType := range ValidNotificationTypes {
		if t == validType {
			return true
		}
	}
	return false
}
