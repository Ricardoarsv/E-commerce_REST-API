package main

import (
	"fmt"
	"log"

	"github.com/Ricardoarsv/E-commerce_REST-API/cmd/api"
	"github.com/Ricardoarsv/E-commerce_REST-API/config"
	"github.com/Ricardoarsv/E-commerce_REST-API/db"
)

func main() {
	db, err := db.NewPostgreSQLStorage(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Envs.DBHost, config.Envs.DBPort, config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBName))

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewApiServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
