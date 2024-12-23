package migrations

import (
	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
)

func Up() {
	db := database.GetDb()

	tables := []interface{}{}
	hunt := model.Hunt{}

	if !db.Migrator().HasTable(hunt) {
		tables = append(tables, hunt)
	}

	db.Migrator().CreateTable(tables...)
	print("migrating table")
}

func Down() {
}
