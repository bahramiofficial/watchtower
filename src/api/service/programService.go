package service

import (
	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"gorm.io/gorm"
)

type ProgramService struct {
	database *gorm.DB
}

func NewProgramService() *ProgramService {
	return &ProgramService{
		database: database.GetDb(),
	}
}

// GetAll retrieves all programs from the database
func (s *ProgramService) GetAllPrograms() ([]model.Program, error) {
	var programs []model.Program
	if err := s.database.Order("id DESC").Find(&programs).Error; err != nil {
		return nil, err
	}
	return programs, nil
}
