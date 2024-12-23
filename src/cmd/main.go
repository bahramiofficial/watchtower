package main

import (
	"log"

	"github.com/bahramiofficial/watchtower/src/api"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/database/migrations"
)

func main() {
	err := database.InitDb()
	if err != nil {
		log.Fatalf("Failed to initialize database")
		print(err)
	}
	//6
	migrations.Up()
	api.InitServer()
}
