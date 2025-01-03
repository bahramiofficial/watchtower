package main

import (
	"log"

	"github.com/bahramiofficial/watchtower/src/api"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/database/migrations"
)

// "github.com/bahramiofficial/watchtower/src/database"
// "github.com/bahramiofficial/watchtower/src/database/migrations"

func main() {
	//watch.EnumAll()
	RunServer()
}

func RunServer() {
	err := database.InitDb()
	if err != nil {
		log.Fatalf("Failed to initialize database")
		print(err)
	}

	//6  6
	migrations.Up()
	api.InitServer()
}
