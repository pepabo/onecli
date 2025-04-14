package utils

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name    string
		fetcher func(page int) ([]models.User, error)
		want    []models.User
		wantErr bool
	}{
		{
			name: "正常系: 1ページ分のデータ",
			fetcher: func(page int) ([]models.User, error) {
				if page == 1 {
					return []models.User{
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
					}, nil
				}
				return []models.User{}, nil
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
			name: "正常系: 複数ページのデータ",
			fetcher: func(page int) ([]models.User, error) {
				if page == 1 {
					return []models.User{
						{
							ID:        1,
							Username:  "testuser1",
							Email:     "test1@example.com",
							Firstname: "Test",
							Lastname:  "User1",
						},
					}, nil
				}
				if page == 2 {
					return []models.User{
						{
							ID:        2,
							Username:  "testuser2",
							Email:     "test2@example.com",
							Firstname: "Test",
							Lastname:  "User2",
						},
					}, nil
				}
				return []models.User{}, nil
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
			name: "異常系: エラー発生",
			fetcher: func(page int) ([]models.User, error) {
				return nil, assert.AnError
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Paginate(tt.fetcher, 1)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
