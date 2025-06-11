package main

import (
	"log"

	"github.com/korolev-n/merch/internal/db"
)

func main() {
	if err := db.SeedInventoryData(); err != nil {
		log.Fatal("seeding failed:", err)
	}
}
