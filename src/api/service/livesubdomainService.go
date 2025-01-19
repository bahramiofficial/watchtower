package service

import (
	"time"

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

type LiveSubdomainFilter struct {
	ProgramName string
	Scope       string
	Provider    string
	Tag         string
	Fresh       bool
	Count       bool
	Limit       int
	Page        int
}

func (s *LiveSubdomainService) GetLiveSubdomains(filter LiveSubdomainFilter) ([]model.LiveSubdomain, int64, error) {
	var subdomains []model.LiveSubdomain
	var count int64

	query := s.database.Model(&model.LiveSubdomain{})

	// Apply filters
	if filter.ProgramName != "" {
		query = query.Where("program_name = ?", filter.ProgramName)
	}
	if filter.Scope != "" {
		query = query.Where("scope = ?", filter.Scope)
	}
	if filter.Provider != "" {

		// Fetch subdomains related to the provider (assuming there's a Subdomains model with a providers field)
		var subdomains []model.Subdomain
		err := s.database.Where("ARRAY_TO_STRING(scopes, ',') ILIKE ?", "%"+filter.Provider+"%").Find(&subdomains).Error
		if err != nil {
			return nil, 0, err
		}

		// Extract the subdomains' names
		subdomainNames := make([]string, len(subdomains))
		for i, subdomain := range subdomains {
			subdomainNames[i] = subdomain.SubDomain
		}

		// Apply the subdomain filter
		query = query.Where("subdomain IN ?", subdomainNames)

		// Apply filter for last_update in the last 12 hours
		twelveHoursAgo := time.Now().Add(-12 * time.Hour)
		query = query.Where("last_update >= ?", twelveHoursAgo)
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

func (s *LiveSubdomainService) GetSingleLiveSubdomainBySubDomain(subdomain string) (model.LiveSubdomain, error) {
	var LiveSubdomain model.LiveSubdomain
	if err := s.database.Where("subdomain = ?", subdomain).First(&LiveSubdomain).Error; err != nil {
		return model.LiveSubdomain{}, err
	}
	return LiveSubdomain, nil
}
