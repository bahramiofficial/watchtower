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

// UpsertHttp performs an upsert (insert or update) operation on the Http model
func UpsertHttp(db *gorm.DB, http Http) {
	var existing Http

	// Check if the record already exists
	err := db.Where("program_name = ? AND sub_domain = ?", http.ProgramName, http.SubDomain).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// New record, create it
			if err := db.Create(&http).Error; err != nil {
				fmt.Printf("Failed to create new Http record: %v", err)
			} else {
				fmt.Printf("New HTTP record created for %s", http.SubDomain)

			}
		} else {
			fmt.Printf("Error finding existing Http record: %v", err)
		}
		return
	}

	// Compare and update fields
	if http.StatusCode != existing.StatusCode {
		fmt.Printf("Change status code for  %s", http.SubDomain)
	}
	if http.Title != existing.Title {
		fmt.Printf("Change title for  %s", http.SubDomain)
	}

	// Perform the update
	if err := db.Model(&existing).Updates(http).Error; err != nil {
		fmt.Printf("Failed to update Http record: %v", err)
		return
	}

}
