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
	LiveSubdomains := model.LiveSubdomain{}
	Http := model.Http{}

	if !db.Migrator().HasTable(Program) {
		tables = append(tables, Program)
	}

	if !db.Migrator().HasTable(Subdomain) {
		tables = append(tables, Subdomain)
	}
	if !db.Migrator().HasTable(LiveSubdomains) {
		tables = append(tables, LiveSubdomains)
	}
	if !db.Migrator().HasTable(Http) {
		tables = append(tables, Http)
	}

	db.Migrator().CreateTable(tables...)
	print("migrating table")
}

func Down() {
}
