package main

import (
	"log"

	"github.com/bahramiofficial/watchtower/src/api"
	"github.com/bahramiofficial/watchtower/src/database"
	"github.com/bahramiofficial/watchtower/src/database/migrations"
	// watchsyncprograms "github.com/bahramiofficial/watchtower/src/utilities/programs"
)

// "github.com/bahramiofficial/watchtower/src/database"
// "github.com/bahramiofficial/watchtower/src/database/migrations"

func main() {
	// watchhttp.Httpx("voorivex.academy")
	// watchsyncprograms.SyncProgramToDb("")
	RunServer()
}

func RunServer() {
	err := database.InitDb()
	if err != nil {
		log.Fatalf("Failed to initialize database")
		print(err)
	}

	//6  6
	// watch.SyncProgramToDb("")
	migrations.Up()
	api.InitServer()
}
