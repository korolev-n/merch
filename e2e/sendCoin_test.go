package e2e

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func registerUser(t *testing.T, username, password string) string {
	body := map[string]string{"username": username, "password": password}
	b, _ := json.Marshal(body)

	resp, err := http.Post(fmt.Sprintf("%s/api/auth", baseURL), "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("register %s failed: %v", username, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("register %s failed with status %d: %s", username, resp.StatusCode, body)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("decode token failed: %v", err)
	}
	return authResp.Token
}

func getBalance(t *testing.T, token string) int {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/info", baseURL), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("get /api/info failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("unexpected status /api/info: %d, body: %s", resp.StatusCode, body)
	}

	var info InfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("decode info response failed: %v", err)
	}
	return info.Coins
}

func deleteUserFromDB(t *testing.T, username string) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not loaded:", err)
	}
	dsn := os.Getenv("MERCH_DB_DSN")
	if dsn == "" {
		t.Fatal("MERCH_DB_DSN not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE username = $1", username)
	if err != nil {
		t.Errorf("failed to delete user %s: %v", username, err)
	}
}

func TestSendCoins(t *testing.T) {
	sender := "sender123"
	receiver := "receiver456"
	password := "testpass"

	//defer deleteUserFromDB(t, sender)
	//defer deleteUserFromDB(t, receiver)

	senderToken := registerUser(t, sender, password)
	receiverToken := registerUser(t, receiver, password)

	initialSenderBalance := getBalance(t, senderToken)
	initialReceiverBalance := getBalance(t, receiverToken)

	sendReq := map[string]interface{}{
		"toUser": receiver,
		"amount": 100,
	}
	sendBytes, _ := json.Marshal(sendReq)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/sendCoin", baseURL), bytes.NewReader(sendBytes))
	if err != nil {
		t.Fatalf("create sendCoin request failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+senderToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("sendCoin request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("sendCoin failed with status %d: %s", resp.StatusCode, body)
	}

	finalSenderBalance := getBalance(t, senderToken)
	finalReceiverBalance := getBalance(t, receiverToken)

	expectedSender := initialSenderBalance - 100
	expectedReceiver := initialReceiverBalance + 100

	if finalSenderBalance != expectedSender {
		t.Errorf("sender balance incorrect: expected %d, got %d", expectedSender, finalSenderBalance)
	}
	if finalReceiverBalance != expectedReceiver {
		t.Errorf("receiver balance incorrect: expected %d, got %d", expectedReceiver, finalReceiverBalance)
	}
}
