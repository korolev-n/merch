package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func SeedInventoryData() error {
	db, err := sql.Open("postgres", os.Getenv("MERCH_DB_DSN"))
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
			"INSERT INTO inventory (type, price) VALUES ($1, $2) ON CONFLICT (inventory_name) DO NOTHING",
			name, price)
		if err != nil {
			log.Printf("Failed to insert %s: %v", name, err)
		} else {
			log.Printf("Inserted: %s - %d", name, price)
		}
	}

	return nil
}
