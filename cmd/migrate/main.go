package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Ricardoarsv/E-commerce_REST-API/config"
	"github.com/Ricardoarsv/E-commerce_REST-API/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// todo implement migration
	db, err := db.NewPostgreSQLStorage(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Envs.DBHost, config.Envs.DBPort, config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBName))

	if err != nil {
		log.Fatal(err)
	}

	// Configurar la migración
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println("Error creating migration instance")

		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]

	switch cmd {
	case "up":
		fmt.Println("Running up migration")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Migration ran successfully!")
	case "down":
		fmt.Println("Running down migration")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Migration ran successfully!")
	default:
		fmt.Println("Invalid command")
	}

}
