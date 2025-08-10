package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"wallet-user-svc/pkg/migrate"
)

func main() {
	var (
		databaseURL    = flag.String("database", "", "Database connection URL (e.g., postgres://user:pass@localhost:5432/dbname?sslmode=disable)")
		migrationsPath = flag.String("path", "./db/migrations", "Path to migrations directory")
		action         = flag.String("action", "up", "Migration action: up, down, status")
		steps          = flag.Int("steps", 1, "Number of migrations to apply/rollback (for down action)")
	)
	flag.Parse()

	if *databaseURL == "" {
		log.Fatal("Database URL is required. Use -database flag")
	}

	config := migrate.Config{
		DatabaseURL:    *databaseURL,
		MigrationsPath: *migrationsPath,
	}

	switch *action {
	case "up":
		if err := migrate.RunMigrations(config); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		if err := migrate.RollbackMigrations(config, *steps); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		fmt.Printf("Rolled back %d migrations successfully\n", *steps)

	case "status":
		version, dirty, err := migrate.GetMigrationStatus(config)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
		if dirty {
			fmt.Printf("Migration version: %d (dirty)\n", version)
		} else {
			fmt.Printf("Migration version: %d\n", version)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s. Use up, down, or status\n", *action)
		os.Exit(1)
	}
}
