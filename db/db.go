package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	t "notification_service/types"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Database defines the interface for database operations
type Database interface {
	GetRateLimitRule(ctx context.Context, nType t.NotificationType) (t.RateLimitRule, error)
	RecordNotification(ctx context.Context, input *t.InputInfo, notif t.Notifications)
	GetLastNotification(ctx context.Context, current t.InputInfo) (t.Notifications, error)
}

type DBConnector struct {
	DB     *sqlx.DB
	Logger *zap.Logger
}

func (db *DBConnector) GetRateLimitRules(ctx context.Context, nType t.NotificationType) (t.RateLimitRule, error) {
	var rule t.RateLimitRule
	query := `
		SELECT  *
		FROM notification_service.rate_limit_rules
		WHERE notification_type = $1
	`

	err := db.DB.GetContext(ctx, &rule, query, nType)
	if err != nil {
		db.Logger.Error("Error fetching rate limit rule",
			zap.Error(err),
			zap.String("notification_type", string(nType)),
		)
		return t.RateLimitRule{}, err
	}
	return rule, err
}

func (db *DBConnector) GetLastNotification(ctx context.Context, current t.InputInfo) (t.Notifications, error) {
	var notif t.Notifications
	query := `
		SELECT *
		FROM notification_service.notifications
		WHERE 
			notification_type = $1 AND
			recipient = $2
		order by created_at desc
		limit 1
	`
	err := db.DB.GetContext(ctx, &notif, query, current.NotificationGroup, current.Recipient)
	if err != nil {
		if err != sql.ErrNoRows {
			db.Logger.Error("Error fetching notifications",
				zap.Error(err),
				zap.String("notification_type", string(current.NotificationGroup)),
			)
			return t.Notifications{}, err
		}
		return t.Notifications{}, nil // not records yet
	}
	db.Logger.Info("Last notification was fetched",
		zap.String("notification_id", notif.ID),
		zap.String("counter", fmt.Sprint(notif.Counter)),
	)
	return notif, err
}

func (db *DBConnector) RecordNotification(ctx context.Context, input *t.InputInfo, notif t.Notifications) error {

	var (
		query string
		err   error
	)
	now := time.Now().UTC()
	if notif.Counter == 0 {
		query = `
		INSERT INTO notification_service.notifications ( recipient, notification_type, counter, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		`
		_, err = db.DB.ExecContext(ctx, query, input.Recipient, input.NotificationGroup, 1, now, now)
	} else {
		query = `
			UPDATE notification_service.notifications 
			SET 
				updated_at = $1, counter = $2
			WHERE
				id = $3
		`
		_, err = db.DB.ExecContext(ctx, query, now, notif.Counter, notif.ID)

	}

	if err != nil {
		db.Logger.Error("Error upserting notification",
			zap.Error(err),
			zap.String("recipient", input.Recipient),
			zap.String("group_type", fmt.Sprint(input.NotificationGroup)),
		)
		return fmt.Errorf("error upserting notification: %w", err)
	}

	db.Logger.Info("Notification upserted successfully",
		zap.String("recipient", input.Recipient),
		zap.String("group_type", fmt.Sprint(input.NotificationGroup)),
	)
	return nil
}
