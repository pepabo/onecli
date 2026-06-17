package onelogin

import (
	"testing"

	"github.com/onelogin/onelogin-go-sdk/v4/pkg/onelogin/models"
	"github.com/stretchr/testify/assert"
)

func TestUserPayload(t *testing.T) {
	user := models.User{
		Firstname: "onecli_test",
		Lastname:  "inatchi",
		Email:     "test@inatchi.dev",
	}

	payload, err := userPayload(user)
	assert.NoError(t, err)

	// The fields explicitly set by the caller must be preserved.
	assert.Equal(t, "onecli_test", payload["firstname"])
	assert.Equal(t, "inatchi", payload["lastname"])
	assert.Equal(t, "test@inatchi.dev", payload["email"])

	// Zero-value time.Time fields must be stripped so OneLogin does not store
	// bogus "0001-01-01T00:00:00Z" timestamps (e.g. last_login).
	for _, k := range []string{
		"created_at",
		"updated_at",
		"activated_at",
		"last_login",
		"password_changed_at",
		"locked_until",
		"invitation_sent_at",
	} {
		_, ok := payload[k]
		assert.Falsef(t, ok, "zero-value time field %q should be stripped from payload", k)
	}
}

func TestUserPayloadKeepsNonZeroTime(t *testing.T) {
	user := models.User{
		Email:  "test@inatchi.dev",
		Status: 2,
	}

	payload, err := userPayload(user)
	assert.NoError(t, err)

	// Non-time fields with real values survive (e.g. set-status sends status only).
	assert.Equal(t, float64(2), payload["status"])
}
