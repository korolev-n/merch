package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func SeedInventoryData() error {
	_ = godotenv.Load()

	dsn := os.Getenv("MERCH_DB_DSN")
	if dsn == "" {
		log.Fatal("MERCH_DB_DSN is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	items := map[string]int{
		"t-shirt":    80,
		"cup":        20,
		"book":       50,
		"pen":        10,
		"powerbank":  200,
		"hoody":      300,
		"umbrella":   200,
		"socks":      10,
		"wallet":     50,
		"pink-hoody": 500,
	}

	for name, price := range items {
		_, err := db.Exec(
			"INSERT INTO inventory (type, price) VALUES ($1, $2) ON CONFLICT (type) DO NOTHING",
			name, price)
		if err != nil {
			log.Printf("Failed to insert %s: %v", name, err)
		}
	}

	return nil
}
