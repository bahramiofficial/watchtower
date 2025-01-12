package model

import (
	"fmt"

	"gorm.io/gorm"
)

// Http represents the HTTP model
type Http struct {
	BaseModel

	ProgramName string      `gorm:"type:text;not null;uniqueIndex:idx_program_subdomain"`
	SubDomain   string      `gorm:"type:text;not null;uniqueIndex:idx_program_subdomain"`
	Scope       string      `gorm:"type:text;not null"`
	IPs         StringArray `gorm:"type:text[]"`
	Tech        StringArray `gorm:"type:text[]"`
	Title       string      `gorm:"type:text"`
	StatusCode  string      `gorm:"type:text"`
	Headers     MapField    `gorm:"type:jsonb"`
	URL         string      `gorm:"type:text"`
	FinalURL    string      `gorm:"type:text"`
	Favicon     string      `gorm:"type:text"`
}

func GetAllHttpWithScope(db *gorm.DB, scope string) ([]Http, error) {
	// Initialize a variable to hold the results
	var https []Http

	// Fetch all https where the scope matches the provided scope
	err := db.Where("scope = ?", scope).Find(&https).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch http for scope '%s': %w", scope, err)
	}

	// Return  Https
	return https, nil
}

func GetAllHttp(db *gorm.DB) ([]Http, error) {
	// Initialize a variable to hold the results
	var https []Http

	// Fetch all https
	err := db.Find(&https).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all http")
	}

	// Return  Https
	return https, nil
}
