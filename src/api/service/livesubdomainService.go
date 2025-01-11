package service

import (
	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"gorm.io/gorm"
)

type LiveSubdomainService struct {
	database *gorm.DB
}

func NewLiveSubdomainService() *LiveSubdomainService {
	return &LiveSubdomainService{
		database: database.GetDb(),
	}
}

// GetLiveSubdomainsByScope retrieves all live subdomains where the scope matches the input
func (s *LiveSubdomainService) GetLiveSubdomainsByScope(scope string) ([]model.LiveSubdomains, error) {
	var liveSubdomains []model.LiveSubdomains
	if err := s.database.Where("scope = ?", scope).Order("id DESC").Find(&liveSubdomains).Error; err != nil {
		return nil, err
	}
	return liveSubdomains, nil
}

// GetLiveSubdomainsByProgramName retrieves all live subdomains where the program name matches the input
func (s *LiveSubdomainService) GetLiveSubdomainsByProgramName(programName string) ([]model.LiveSubdomains, error) {
	var liveSubdomains []model.LiveSubdomains
	if err := s.database.Where("program_name = ?", programName).Order("id DESC").Find(&liveSubdomains).Error; err != nil {
		return nil, err
	}
	return liveSubdomains, nil
}
