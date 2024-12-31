package migrations

import (
	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
)

func Up() {
	db := database.GetDb()

	tables := []interface{}{}
	Program := model.Program{}
	Subdomain := model.Subdomain{}

	if !db.Migrator().HasTable(Program) {
		tables = append(tables, Program)
	}

	if !db.Migrator().HasTable(Subdomain) {
		tables = append(tables, Subdomain)
	}

	db.Migrator().CreateTable(tables...)
	print("migrating table")
}

func Down() {
}
