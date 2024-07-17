package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"notification_service/server"
	"notification_service/types"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNotificationService struct {
	mock.Mock
}

// SendNotification simula el método SendNotification.
func (m *MockNotificationService) SendNotification(ctx context.Context, input types.InputInfo) (types.Output, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(types.Output), args.Error(1)
}

// IsAllowed simula el método IsAllowed.
func (m *MockNotificationService) IsAllowed(ctx context.Context, nType types.NotificationType) bool {
	args := m.Called(ctx, nType)
	return args.Bool(0)
}

type MockDBConnector struct {
	mock.Mock
}

func (m *MockDBConnector) Query(query string, args ...interface{}) (interface{}, error) {
	argsMock := m.Called(query, args)
	return argsMock.Get(0), argsMock.Error(1)
}

func (m *MockDBConnector) Close() error {
	argsMock := m.Called()
	return argsMock.Error(0)
}
func TestValidateInputData(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedError  string // Expected error substring
		expectedFields types.InputInfo
	}{
		{
			name:           "Valid Input",
			input:          `{"recipient": "example@example.com", "group": "Marketing"}`,
			expectedError:  "", // Happy path
			expectedFields: types.InputInfo{Recipient: "example@example.com", NotificationGroup: "Marketing"},
		},
		{
			name:          "Invalid JSON",
			input:         `invalid json`,
			expectedError: "invalid JSON format",
		},
		{
			name:          "Missing Required Fields",
			input:         `{"recipient": "", "group": ""}`,
			expectedError: "missing required fields",
		},
		{
			name:          "Invalid Notification Type",
			input:         `{"recipient": "example@example.com", "group": "invalidType"}`,
			expectedError: "invalid notification type: invalidType",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.NewReader(tc.input)
			in, err := server.ValidateInputData(body)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedFields, in)
			}
		})
	}
}
func TestSendHandler(t *testing.T) {
	mockService := &MockNotificationService{}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/notify":
			if r.Method != http.MethodPost {
				http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
				return
			}
			in, err := server.ValidateInputData(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			mockService.On("SendNotification", context.Background(), in).Return(
				types.Output{
					Recipient:         "example@example.com",
					NotificationGroup: "NEWS",
				},
				nil)
			out, err := mockService.SendNotification(context.Background(), in)
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
			if _, err := w.Write(response); err != nil {
				t.Logf("could not send status OK: %v", err)
			}
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}))
	defer s.Close()

	testCases := []struct {
		name             string
		requestBody      string
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "Valid Request",
			requestBody:      `{"recipient": "example@example.com", "group": "Marketing"}`,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"recipient":"example@example.com","group":"NEWS","timestamp":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `invalid json`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Missing Required Fields",
			requestBody:    `{"recipient": "", "group": ""}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Invalid Notification Type",
			requestBody:    `{"recipient": "example@example.com", "group": "invalidType"}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqBody := strings.NewReader(tc.requestBody)
			req, err := http.NewRequest("POST", s.URL+"/notify", reqBody)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %v, got %v", tc.expectedStatus, resp.Status)
			}

			if tc.expectedResponse != "" {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("Failed to read response body: %v", err)
				}
				if string(body) != tc.expectedResponse {
					t.Errorf("Expected response body %q, got %q", tc.expectedResponse, string(body))
				}
			}
		})
	}
}
