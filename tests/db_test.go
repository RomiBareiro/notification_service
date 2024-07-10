package tests

import (
	"context"
	"database/sql/driver"
	d "notification_service/db"
	"notification_service/types"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func TestRecordNotification(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	logger := zap.NewNop()

	connector := &d.DBConnector{
		DB:     db,
		Logger: logger,
	}

	testCases := []struct {
		description string
		counter     int
		expectQuery string
		expectArgs  []driver.Value
	}{
		{
			description: "Counter = 0 (INSERT)",
			counter:     0,
			expectQuery: `INSERT INTO notification_service.notifications`,
		},
		{
			description: "Counter != 0 (UPDATE)",
			counter:     1,
			expectQuery: `UPDATE notification_service.notifications`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mock.ExpectExec(tc.expectQuery).
				WithArgs(tc.expectArgs...).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err := connector.RecordNotification(context.Background(), &types.InputInfo{
				Recipient:         "recipient",
				NotificationGroup: "Marketing",
			}, types.Notifications{
				ID:               "id",
				NotificationType: types.Marketing,
				Counter:          tc.counter,
				Recipient:        "recipient",
			})
			if err != nil {
				t.Fatalf("Error recording notification: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}
