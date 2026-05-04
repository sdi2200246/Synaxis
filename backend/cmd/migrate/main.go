package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	if len(os.Args) < 3 {
    	log.Fatalf("usage: migrate [up|down] [YOUR DB ENV URL]")
	}

	migration := os.Args[1];
	dbEnvVariable := os.Args[2];
	godotenv.Load() 
	dbURL := os.Getenv(dbEnvVariable)
	if dbURL == "" {
		log.Fatalf("Failed to fetch DATABASE_URL from the .env file",)
	}

	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch migration {
		case "up":
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("failed to run up migration: %v", err)
			}
		case "down":
			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("failed to run down migration: %v", err)
			}
		default:
			log.Fatalf("unknown command %q, use 'up' or 'down'", migration)
		}

	log.Println("migrations applied successfully")
}