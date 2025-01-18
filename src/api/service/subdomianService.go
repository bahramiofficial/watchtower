package service

import (
	"time"

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

func (s *SubdomainService) GetSingleSubdomainBySubDomain(subdomain string) (model.Subdomain, error) {
	var Subdomain model.Subdomain
	if err := s.database.Where("subdomain = ?", subdomain).First(&Subdomain).Error; err != nil {
		return model.Subdomain{}, err
	}
	return Subdomain, nil
}

type SubdomainFilter struct {
	ProgramName string
	Scope       string
	Provider    string
	Fresh       bool
	Count       bool
	Limit       int
	Page        int
}

func (s *SubdomainService) GetSubdomains(filter SubdomainFilter) ([]model.Subdomain, int64, error) {
	var subdomains []model.Subdomain
	var count int64

	query := s.database.Model(&model.Subdomain{})

	// Apply filters
	if filter.ProgramName != "" {
		query = query.Where("program_name = ?", filter.ProgramName)
	}
	if filter.Scope != "" {
		query = query.Where("scope = ?", filter.Scope)
	}
	if filter.Provider != "" {
		query = query.Where("? = ANY(providers)", filter.Provider) // Use PostgreSQL's ANY() for array search
	}
	if filter.Fresh {
		twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)
		query = query.Where("created_at >= ?", twentyFourHoursAgo)
	}

	// Get count if requested
	if filter.Count {
		if err := query.Count(&count).Error; err != nil {
			return nil, 0, err
		}
		return nil, count, nil
	}

	// Apply pagination if specified
	if filter.Limit > 0 && filter.Page > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Execute query
	if err := query.Find(&subdomains).Error; err != nil {
		return nil, 0, err
	}

	return subdomains, int64(len(subdomains)), nil
}
