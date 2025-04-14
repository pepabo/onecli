package utils

import (
	"bytes"
	"testing"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestPrintOutput(t *testing.T) {
	tests := []struct {
		name     string
		format   OutputFormat
		data     any
		expected any
		wantErr  bool
	}{
		{
			name:   "正常系: JSON形式で出力",
			format: OutputFormatJSON,
			data: []models.User{
				{
					ID:                1,
					Username:          "testuser1",
					Email:             "test1@example.com",
					Firstname:         "Test",
					Lastname:          "User1",
					CreatedAt:         time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt:         time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC),
					ActivatedAt:       time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC),
					LastLogin:         time.Time{},
					PasswordChangedAt: time.Time{},
					LockedUntil:       time.Time{},
					InvitationSentAt:  time.Time{},
					State:             1,
					Status:            1,
				},
			},
			expected: `[
  {
    "id": 1,
    "username": "testuser1",
    "email": "test1@example.com",
    "firstname": "Test",
    "lastname": "User1",
    "created_at": "2024-04-01T12:00:00Z",
    "updated_at": "2024-04-01T12:00:00Z",
    "activated_at": "2024-04-01T12:00:00Z",
    "last_login": "0001-01-01T00:00:00Z",
    "password_changed_at": "0001-01-01T00:00:00Z",
    "locked_until": "0001-01-01T00:00:00Z",
    "invitation_sent_at": "0001-01-01T00:00:00Z",
    "state": 1,
    "status": 1
  }
]`,
			wantErr: false,
		},
		{
			name:   "正常系: YAML形式で出力",
			format: OutputFormatYAML,
			data: []models.User{
				{
					ID:          1,
					Username:    "testuser1",
					Email:       "test1@example.com",
					Firstname:   "Test",
					Lastname:    "User1",
					CreatedAt:   time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC),
					ActivatedAt: time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC),
					State:       1,
					Status:      1,
				},
			},
			expected: []map[string]interface{}{
				{
					"id":           1,
					"username":     "testuser1",
					"email":        "test1@example.com",
					"firstname":    "Test",
					"lastname":     "User1",
					"created_at":   time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC),
					"updated_at":   time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC),
					"activated_at": time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC),
					"state":        1,
					"status":       1,
				},
			},
			wantErr: false,
		},
		{
			name:     "正常系: 空のデータ",
			format:   OutputFormatJSON,
			data:     []models.User{},
			expected: "[]",
			wantErr:  false,
		},
		{
			name:     "異常系: 不正なデータ",
			format:   OutputFormatJSON,
			data:     make(chan int), // エンコードできない型
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := PrintOutput(tt.data, tt.format, &buf)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			output := buf.String()

			if tt.format == OutputFormatJSON {
				assert.JSONEq(t, tt.expected.(string), output)
			} else {
				expectedBytes, err := yaml.Marshal(tt.expected)
				assert.NoError(t, err)
				assert.YAMLEq(t, string(expectedBytes), output)
			}
		})
	}
}
