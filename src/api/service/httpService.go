package service

import (
	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"gorm.io/gorm"
)

type HttpService struct {
	database *gorm.DB
}

func NewHttpService() *HttpService {
	return &HttpService{
		database: database.GetDb(),
	}
}

// GetHttpByScope retrieves all HTTP data where the scope matches the input
func (s *HttpService) GetHttpByScope(scope string) ([]model.Http, error) {
	var httpData []model.Http
	if err := s.database.Where("scope = ?", scope).Order("id DESC").Find(&httpData).Error; err != nil {
		return nil, err
	}
	return httpData, nil
}

// GetHttpByProgramName retrieves all HTTP data where the program name matches the input
func (s *HttpService) GetHttpByProgramName(programName string) ([]model.Http, error) {
	var httpData []model.Http
	if err := s.database.Where("program_name = ?", programName).Order("id DESC").Find(&httpData).Error; err != nil {
		return nil, err
	}
	return httpData, nil
}
