package service

import (
	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/database"
	"gorm.io/gorm"
)

type SubdomainService struct {
	database *gorm.DB
}

func NewSubdomainService() *SubdomainService {
	return &SubdomainService{
		database: database.GetDb(),
	}
}

// GetSubdomainsByScope retrieves all subdomains where the scope matches the input
func (s *SubdomainService) GetSubdomainsByScope(scope string) ([]model.Subdomain, error) {
	var subdomains []model.Subdomain
	if err := s.database.Where("scope = ?", scope).Order("id DESC").Find(&subdomains).Error; err != nil {
		return nil, err
	}
	return subdomains, nil
}

func (s *SubdomainService) GetSubdomainsByProgramName(programName string) ([]model.Subdomain, error) {
	var subdomains []model.Subdomain
	if err := s.database.Where("programname = ?", programName).Order("id DESC").Find(&subdomains).Error; err != nil {
		return nil, err
	}
	return subdomains, nil
}

func (s *SubdomainService) GetAllSubdomain() ([]model.Subdomain, error) {
	var subdomains []model.Subdomain
	if err := s.database.Order("id DESC").Find(&subdomains).Error; err != nil {
		return nil, err
	}
	return subdomains, nil
}