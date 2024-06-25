package main

import (
	"github.com/Karol7Krawczyk/golang-migrate/migrations/config"
	"github.com/Karol7Krawczyk/golang-migrate/migrations/db"
	"github.com/Karol7Krawczyk/golang-migrate/migrations/handlers"

	"log"
)

func main() {
	config := config.ParseFlags()

	database, err := db.GetConnection(config)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.CloseConnection(database)

	if err := db.PrepareMigrationTable(database, config); err != nil {
		log.Fatalf("Error preparing migration table: %v", err)
	}

	handlers.HandleCommand(database, config)
}
