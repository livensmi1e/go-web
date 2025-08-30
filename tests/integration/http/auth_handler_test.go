package http_test

import (
	"go-web/tests/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthFlow(t *testing.T) {
	ts := utils.SetupTestServer()
	defer ts.Server.Close()

	email := utils.GenUserEmail()
	password := "password123"
	t.Run("should register a new user successfully", func(t *testing.T) {
		registerBody := map[string]string{
			"email":    email,
			"password": password,
		}
		var registerResp map[string]any
		ts.DoRequest(t, "POST", "/api/auth/register", registerBody, "", &registerResp, 201)
		data := registerResp["data"].(map[string]any)
		require.NotEmpty(t, data["id"])
		require.Equal(t, registerBody["email"], data["email"])
	})

	var token string
	t.Run("should login sucessfully", func(t *testing.T) {
		loginBody := map[string]string{
			"email":    email,
			"password": password,
		}
		var loginResp map[string]any
		ts.DoRequest(t, "POST", "/api/auth/login", loginBody, "", &loginResp, 200)
		token = loginResp["data"].(map[string]any)["token"].(string)
		require.NotEmpty(t, token)
	})

	t.Run("should return get me when have token", func(t *testing.T) {
		var meResp map[string]any
		ts.DoRequest(t, "GET", "/api/me", nil, token, &meResp, 200)
		require.Equal(t, "me", meResp["data"])
	})

	// Need some tests for other cases like invalid input, wrong password, no token, invalid token, rate limit, etc.
}
