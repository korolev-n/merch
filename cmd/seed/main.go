package main

import (
	"log"

	"github.com/korolev-n/merch-auth/internal/db"
)

func main() {
	if err := db.SeedInventoryData(); err != nil {
		log.Fatal("seeding failed:", err)
	}
}
