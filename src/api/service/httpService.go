package service

import (
	"time"

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

type HttpFilter struct {
	ProgramName string
	Scope       string
	Provider    string
	Title       string
	Status      string
	Tech        string
	Latest      bool
	Fresh       bool
	Count       bool
	Limit       int
	Page        int
}

func (s *HttpService) GetHttp(filter HttpFilter) ([]model.Http, int64, error) {
	var https []model.Http
	var count int64

	// Initialize query
	query := s.database.Model(&model.Http{})

	// Apply filters
	if filter.ProgramName != "" {
		query = query.Where("program_name = ?", filter.ProgramName)
	}
	if filter.Scope != "" {
		query = query.Where("scope = ?", filter.Scope)
	}
	if filter.Title != "" {
		query = query.Where("title = ?", filter.Title)
	}
	if filter.Status != "" {
		query = query.Where("status_code = ?", filter.Status)
	}
	if filter.Tech != "" {
		query = query.Where("? = ANY(tech)", filter.Tech) // Correct handling of tech array
	}
	if filter.Fresh {
		twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)
		query = query.Where("created_at >= ?", twentyFourHoursAgo) // Use `created_at` assuming a standard timestamp field
	}
	if filter.Latest {
		twelveHoursAgo := time.Now().Add(-12 * time.Hour)
		query = query.Where("updated_at >= ?", twelveHoursAgo) // Use `updated_at` for last updates
	}
	if filter.Provider != "" {
		// Fetch subdomains related to the provider
		var subdomains []model.Subdomain
		err := s.database.Where("ARRAY_TO_STRING(providers, ',') ILIKE ?", "%"+filter.Provider+"%").Find(&subdomains).Error
		if err != nil {
			return nil, 0, err
		}

		// Extract subdomain names
		subdomainNames := make([]string, len(subdomains))
		for i, subdomain := range subdomains {
			subdomainNames[i] = subdomain.SubDomain
		}

		if len(subdomainNames) > 0 {
			query = query.Where("sub_domain IN ?", subdomainNames) // Use `sub_domain` to match your struct
		}
	}

	// Count records if requested
	if filter.Count {
		err := query.Count(&count).Error
		if err != nil {
			return nil, 0, err
		}
		return nil, count, nil
	}

	// Apply pagination
	if filter.Limit > 0 && filter.Page > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Fetch results
	err := query.Find(&https).Error
	if err != nil {
		return nil, 0, err
	}

	return https, count, nil
}
func (s *HttpService) GetSingleHttpBySubDomain(subdomain string) (model.Http, error) {
	var httpModel = model.Http{}
	if err := s.database.Where("subdomain = ?", subdomain).First(&httpModel).Error; err != nil {
		return model.Http{}, err
	}
	return httpModel, nil
}
