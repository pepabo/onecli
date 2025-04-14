package utils

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestConvertToUsers(t *testing.T) {
	tests := []struct {
		name    string
		input   []interface{}
		want    []models.User
		wantErr bool
	}{
		{
			name: "正常系: ユーザー情報の変換",
			input: []interface{}{
				map[string]interface{}{
					"id":        float64(1),
					"username":  "testuser1",
					"email":     "test1@example.com",
					"firstname": "Test",
					"lastname":  "User1",
				},
				map[string]interface{}{
					"id":        float64(2),
					"username":  "testuser2",
					"email":     "test2@example.com",
					"firstname": "Test",
					"lastname":  "User2",
				},
			},
			want: []models.User{
				{
					ID:        1,
					Username:  "testuser1",
					Email:     "test1@example.com",
					Firstname: "Test",
					Lastname:  "User1",
				},
				{
					ID:        2,
					Username:  "testuser2",
					Email:     "test2@example.com",
					Firstname: "Test",
					Lastname:  "User2",
				},
			},
			wantErr: false,
		},
		{
			name: "異常系: 不正なデータ",
			input: []interface{}{
				"invalid data",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToUsers(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
