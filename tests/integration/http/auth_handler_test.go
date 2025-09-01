package http_test

import (
	"testing"
	"time"

	"go-web/tests/utils"

	"github.com/stretchr/testify/require"
)

func TestAuthFlow(t *testing.T) {
	ts := utils.SetupTestServer()
	defer ts.Server.Close()

	email := utils.GenUserEmail()
	password := "password123"
	t.Run("register a new user successfully", func(t *testing.T) {
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
	t.Run("login successfully sets refresh token cookie", func(t *testing.T) {
		loginBody := map[string]string{
			"email":    email,
			"password": password,
		}
		var loginResp map[string]any
		res := ts.DoRequest(t, "POST", "/api/auth/login", loginBody, "", &loginResp, 200)
		token = loginResp["data"].(map[string]any)["token"].(string)
		require.NotEmpty(t, token)
		var hasRefresh bool
		for _, c := range res.Cookies() {
			if c.Name == "refreshToken" {
				hasRefresh = true
				require.NotEmpty(t, c.Value)
			}
		}
		require.True(t, hasRefresh, "refreshToken cookie must be set")
	})

	t.Run("return get me when have token", func(t *testing.T) {
		var meResp map[string]any
		ts.DoRequest(t, "GET", "/api/me", nil, token, &meResp, 200)
		require.Equal(t, "me", meResp["data"])
	})

	t.Run("refresh returns new access token", func(t *testing.T) {
		loginBody := map[string]string{
			"email":    email,
			"password": password,
		}
		var loginResp map[string]any
		res := ts.DoRequest(t, "POST", "/api/auth/login", loginBody, "", &loginResp, 200)
		var refreshResp map[string]any
		res = ts.DoRequest(t, "POST", "/api/auth/refresh", nil, "", &refreshResp, 200, res.Cookies()...)
		newToken := refreshResp["data"].(map[string]any)["token"].(string)
		require.NotEmpty(t, newToken)
		require.NotEqual(t, token, newToken)
		var hasNewRefresh bool
		for _, c := range res.Cookies() {
			if c.Name == "refreshToken" {
				hasNewRefresh = true
				require.NotEmpty(t, c.Value)
			}
		}
		require.True(t, hasNewRefresh, "refreshToken cookie must be rotated")
	})

	t.Run("logout clears refresh token", func(t *testing.T) {
		loginBody := map[string]string{
			"email":    email,
			"password": password,
		}
		var loginResp map[string]any
		res := ts.DoRequest(t, "POST", "/api/auth/login", loginBody, "", &loginResp, 200)
		token := loginResp["data"].(map[string]any)["token"].(string)
		var logoutResp map[string]any
		res = ts.DoRequest(t, "POST", "/api/auth/logout", nil, token, &logoutResp, 200, res.Cookies()...)
		var cleared bool
		for _, c := range res.Cookies() {
			if c.Name == "refreshToken" && c.Value == "" && c.MaxAge == -1 {
				cleared = true
			}
		}
		require.True(t, cleared, "refreshToken cookie must be cleared on logout")
	})

	t.Run("register with invalid email", func(t *testing.T) {
		registerBody := map[string]string{
			"email":    "not-an-email",
			"password": password,
		}
		var resp map[string]any
		ts.DoRequest(t, "POST", "/api/auth/register", registerBody, "", &resp, 400)
	})

	t.Run("register with duplicate email", func(t *testing.T) {
		registerBody := map[string]string{
			"email":    email,
			"password": password,
		}
		var resp map[string]any
		ts.DoRequest(t, "POST", "/api/auth/register", registerBody, "", &resp, 409)
	})

	t.Run("login with wrong password", func(t *testing.T) {
		loginBody := map[string]string{
			"email":    email,
			"password": "wrongpassword",
		}
		var resp map[string]any
		ts.DoRequest(t, "POST", "/api/auth/login", loginBody, "", &resp, 401)
	})

	t.Run("login with non-existent email", func(t *testing.T) {
		loginBody := map[string]string{
			"email":    "unknown@example.com",
			"password": password,
		}
		var resp map[string]any
		ts.DoRequest(t, "POST", "/api/auth/login", loginBody, "", &resp, 401)
	})

	t.Run("get me without token", func(t *testing.T) {
		var resp map[string]any
		ts.DoRequest(t, "GET", "/api/me", nil, "", &resp, 401)
	})

	t.Run("get me with invalid token", func(t *testing.T) {
		var resp map[string]any
		ts.DoRequest(t, "GET", "/api/me", nil, "fake.token.here", &resp, 401)
	})

	t.Run("rate limit exceeded", func(t *testing.T) {
		time.Sleep(3 * time.Second) // wait for previous tests to not affect this one
		for i := 0; i < 30; i++ {
			ts.DoRequest(t, "GET", "/api/me", nil, token, nil, 200)
		}
		var resp map[string]any
		ts.DoRequest(t, "GET", "/api/me", nil, token, &resp, 429)
	})
}
