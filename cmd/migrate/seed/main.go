package main

import (
	"log"

	"github.com/batugkocak/social/internal/db"
	"github.com/batugkocak/social/internal/env"
	"github.com/batugkocak/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://postgres:postgres@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStorage(conn)
	db.Seed(store)
}
