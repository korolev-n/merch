package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestBuyItem(t *testing.T) {
	username := "e2euser"
	password := "testpass"

	// 1. Регистрация пользователя
	registerBody := map[string]string{
		"username": username,
		"password": password,
	}
	registerBytes, _ := json.Marshal(registerBody)
	resp, err := http.Post(fmt.Sprintf("%s/api/auth", baseURL), "application/json", bytes.NewReader(registerBytes))
	if err != nil {
		t.Fatalf("Failed to register: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Unexpected status on register: %d, body: %s", resp.StatusCode, body)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("Failed to decode auth response: %v", err)
	}
	token := authResp.Token
	if token == "" {
		t.Fatal("Token not received")
	}

	// 2. Покупка товара (например, "cup")
	item := "cup"
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/buy/%s", baseURL, item), nil)
	if err != nil {
		t.Fatalf("Failed to create buy request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Buy request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Unexpected status on buy: %d, body: %s", resp.StatusCode, body)
	}

	// 3. Проверка наличия товара в инвентаре
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/info", baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("GetInfo failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to fetch /api/info: %d", resp.StatusCode)
	}

	var info InfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("Failed to parse /api/info response: %v", err)
	}

	found := false
	for _, item := range info.Inventory {
		if item.Type == "cup" && item.Quantity >= 1 {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected item 'cup' in inventory but not found")
	}
}
