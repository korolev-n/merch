package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	userCount    = 100 // целевая 100_000
	defaultPass  = "Password123"
	passwordCost = bcrypt.DefaultCost
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not loaded:", err)
	}

	fmt.Println("MERCH_DB_DSN:", os.Getenv("MERCH_DB_DSN"))

	db, err := sql.Open("postgres", os.Getenv("MERCH_DB_DSN"))
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal("Could not start transaction:", err)
	}

	// Подготовка запросов
	stmtUser, err := tx.PrepareContext(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id")
	if err != nil {
		log.Fatal("User prepare failed:", err)
	}
	defer stmtUser.Close()

	passHash, _ := bcrypt.GenerateFromPassword([]byte(defaultPass), passwordCost)

	// Срезы для пакетной вставки
	var usernames []string
	var passwordHashes []string
	var userIDs []int

	for i := 0; i < userCount; i++ {
		username := fmt.Sprintf("user_%d", i)
		usernames = append(usernames, username)
		passwordHashes = append(passwordHashes, string(passHash))
	}

	// Пакетная вставка пользователей
	for i := 0; i < userCount; i++ {
		var userID int
		err := stmtUser.QueryRowContext(ctx, usernames[i], passwordHashes[i]).Scan(&userID)
		if err != nil {
			log.Fatalf("Insert user failed at %d: %v", i, err)
		}
		userIDs = append(userIDs, userID)

		if i%1000 == 0 {
			fmt.Printf("Inserted %d users\n", i)
		}
	}

	// Подготовка запроса для вставки кошельков
	stmtWallet, err := tx.PrepareContext(ctx, "INSERT INTO wallets (user_id, balance) VALUES ($1, 1000)")
	if err != nil {
		log.Fatal("Wallet prepare failed:", err)
	}
	defer stmtWallet.Close()

	// Пакетная вставка кошельков
	for _, userID := range userIDs {
		_, err = stmtWallet.ExecContext(ctx, userID)
		if err != nil {
			log.Fatalf("Insert wallet failed for user_id %d: %v", userID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Commit failed:", err)
	}
	fmt.Println("Inserted all users and wallets")
}
