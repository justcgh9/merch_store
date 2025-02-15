package e2e__test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

// TestEndToEndFlow covers the happy path and a few negative scenarios
func TestEndToEndFlow(t *testing.T) {
	// Step 1: Register and authenticate a test user
	token := authenticateUser(t, "testuser", "testpass")
	require.NotEmpty(t, token)

	// Step 2: Fetch initial info (coins, inventory, history)
	initialInfo := makeRequest(t, "GET", "/api/info", nil, token, http.StatusOK)
	require.NotNil(t, initialInfo["coins"])

	initialCoins := int(initialInfo["coins"].(float64))

	// Step 3: Send coins to another user (happy path)
	_ = authenticateUser(t, "anotheruser", "password")
	sendCoins(t, token, "anotheruser", 10)

	// Step 4: Validate updated coins for sender
	updatedInfo := makeRequest(t, "GET", "/api/info", nil, token, http.StatusOK)
	updatedCoins := int(updatedInfo["coins"].(float64))
	require.Equal(t, initialCoins-10, updatedCoins)

	// Step 5: Buy an item (happy path)
	makeRequest(t, "GET", "/api/buy/pink-hoody", nil, token, http.StatusOK)

	// ------ Additional Scenarios ------

	// (Negative) Send coins to yourself
	errorResp := makeRequest(t, "POST", "/api/sendCoin", map[string]interface{}{
		"toUser": "testuser",
		"amount": 5,
	}, token, http.StatusBadRequest)
	require.Contains(t, errorResp["errors"], "cannot send money to yourself")

	// (Negative) Send negative amount
	errorResp = makeRequest(t, "POST", "/api/sendCoin", map[string]interface{}{
		"toUser": "anotheruser",
		"amount": -10,
	}, token, http.StatusBadRequest)
	require.Contains(t, errorResp["errors"], "error cannot send less than 0 to another user")

	// (Negative) Buy an item with insufficient coins
	errorResp = makeRequest(t, "GET", "/api/buy/pink-hoody", nil, token, http.StatusBadRequest)
	require.Contains(t, errorResp["errors"], "could not buy pink-hoody")
	
	// (Negative) Attempt buying non-existing item
	errorResp = makeRequest(t, "GET", "/api/buy/non_existing_item", nil, token, http.StatusBadRequest)
	require.Contains(t, errorResp["errors"], "could not buy non_existing_item")

	// (Negative) Unauthorized request (no token)
	makeRequest(t, "GET", "/api/info", nil, "", http.StatusUnauthorized)

	// (Negative) Invalid token
	makeRequest(t, "GET", "/api/info", nil, "invalid.token.here", http.StatusUnauthorized)
}

func authenticateUser(t *testing.T, username, password string) string {
	resp := makeRequest(t, "POST", "/api/auth", map[string]string{
		"username": username,
		"password": password,
	}, "", http.StatusOK)
	token, ok := resp["token"].(string)
	fmt.Println(token)
	require.True(t, ok)
	return token
}

func sendCoins(t *testing.T, token, toUser string, amount int) {
	makeRequest(t, "POST", "/api/sendCoin", map[string]interface{}{
		"toUser": toUser,
		"amount": amount,
	}, token, http.StatusOK)
}

func makeRequest(t *testing.T, method, path string, payload interface{}, token string, expectedStatus int) map[string]interface{} {
	var reqBody []byte
	var err error

	if payload != nil {
		reqBody, err = json.Marshal(payload)
		require.NoError(t, err)
	}

	req, err := http.NewRequest(method, baseURL+path, bytes.NewBuffer(reqBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}

	if (resp.StatusCode != 200) {
		_ = json.NewDecoder(resp.Body).Decode(&result)
		fmt.Println(result["errors"])
	}

	require.Equal(t, expectedStatus, resp.StatusCode)

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest {
		_ = json.NewDecoder(resp.Body).Decode(&result)
	}

	return result
}
