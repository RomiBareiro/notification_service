package tests

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	d "notification_service/db"
	"notification_service/service"
	"notification_service/types"
)

func TestIsAllowed(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")

	logger := zap.NewNop()
	now := time.Now()
	testCases := []struct {
		name              string
		inputType         types.NotificationType
		mockGetLastResult *sqlmock.Rows
		mockGetRateResult *sqlmock.Rows
		expectedAllowed   bool
		expectedError     bool
	}{
		{
			name:              "First Notification",
			inputType:         types.Status,
			mockGetLastResult: sqlmock.NewRows([]string{"id", "notification_type", "counter", "recipient", "created_at", "updated_at"}), // Simulate no rows returned
			mockGetRateResult: sqlmock.NewRows([]string{"id", "notification_type", "max_count", "duration"}).
				AddRow("cbd060ef-f620-489d-8ffc-49f73cde4f54", "Status", 10, 3600.00),
			expectedAllowed: true,
			expectedError:   false,
		},
		{
			name:      "Exceeds Max Count",
			inputType: types.Marketing,
			mockGetLastResult: sqlmock.NewRows([]string{"id", "notification_type", "counter", "recipient", "created_at", "updated_at"}).
				AddRow("cbd060ef-f620-489d-8ffc-49f73cde4f54", "Marketing", 3, "recipient", now.Add(-time.Minute), now),
			mockGetRateResult: sqlmock.NewRows([]string{"id", "notification_type", "max_count", "duration"}).
				AddRow("cbd060ef-f620-489d-8ffc-49f73cde4f54", "Marketing", 3, 100.00),
			expectedAllowed: false,
			expectedError:   false,
		},
		{
			name:      "Error Getting Rate Limit Rules",
			inputType: types.Status,
			mockGetLastResult: sqlmock.NewRows([]string{"id", "notification_type", "counter", "recipient", "created_at", "updated_at"}).
				AddRow("cbd060ef-f620-489d-8ffc-49f73cde4f54", "Status", 0, "recipient", now, now),
			mockGetRateResult: nil,
			expectedAllowed:   false,
			expectedError:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log("Running test case:", tc.name)
			mock.ExpectQuery(`SELECT \* FROM notification_service.notifications`).
				WithArgs(tc.inputType, "recipient").
				WillReturnRows(tc.mockGetLastResult)

			if tc.mockGetRateResult != nil {
				mock.ExpectQuery(`SELECT \* FROM notification_service.rate_limit_rules`).
					WillReturnRows(tc.mockGetRateResult)
			}

			service := &service.NotificationService{
				DB:     &d.DBConnector{DB: db, Logger: logger},
				Logger: logger,
				Input:  types.InputInfo{Recipient: "recipient", NotificationGroup: tc.inputType},
			}

			allowed := service.IsAllowed(context.Background(), tc.inputType)

			if tc.expectedError {
				assert.False(t, allowed, "Expected not allowed due to error")
			} else {
				assert.Equal(t, tc.expectedAllowed, allowed, "Expected allowed status mismatch")
			}

		})
	}
}
