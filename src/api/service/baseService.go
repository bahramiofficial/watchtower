package service

import (
	"github.com/bahramiofficial/watchtower/src/database"
	"gorm.io/gorm"
)

type BaseService struct {
	database *gorm.DB
}

func NewBaseService() *BaseService {
	return &BaseService{
		database: database.GetDb(),
	}
}
