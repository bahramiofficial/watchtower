package model

import (
	"fmt"

	"github.com/bahramiofficial/watchtower/src/utilities"
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
	err := db.Where("sub_domain = ?", http.SubDomain).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {

			// New record, create it
			if err := db.Create(&http).Error; err != nil {
				fmt.Printf("Failed to create new Http record: %v", err)
			} else {
				utilities.SendDiscordMessage(fmt.Sprintf("```'%s' (fresh http) has been added to '%s' program```", http.SubDomain, http.ProgramName))

				fmt.Printf("New HTTP record created for %s", http.SubDomain)

			}
		} else {
			fmt.Printf("Error finding existing Http record: %v", err)
		}
		return
	}

	// Prepare a map to track changes
	updatedFields := map[string]interface{}{}

	// Compare fields and track updates
	if http.StatusCode != existing.StatusCode {
		updatedFields["status_code"] = http.StatusCode
		utilities.SendDiscordMessage(fmt.Sprintf("```'%s' (fresh http) has been updated status code to : '%s'```", http.SubDomain, http.StatusCode))

	}
	if http.Title != existing.Title {
		updatedFields["title"] = http.Title
		utilities.SendDiscordMessage(fmt.Sprintf("```'%s' (fresh http) has been updated title to : '%s'```", http.SubDomain, http.Title))

	}
	if http.Favicon != existing.Favicon {
		updatedFields["favicon"] = http.Favicon
		utilities.SendDiscordMessage(fmt.Sprintf("```'%s' (fresh http) has been updated favicon to : '%s'```", http.SubDomain, http.Favicon))

	}
	if !compareStringSlices(http.IPs, existing.IPs) {
		updatedFields["ips"] = http.IPs
	}
	if !compareStringSlices(http.Tech, existing.Tech) {
		updatedFields["tech"] = http.Tech
	}
	if http.URL != existing.URL {
		updatedFields["url"] = http.URL
	}
	if http.FinalURL != existing.FinalURL {
		updatedFields["final_url"] = http.FinalURL
	}
	if !compareMapFields(http.Headers, existing.Headers) {
		updatedFields["headers"] = http.Headers
	}

	// Only update if there are changes
	if len(updatedFields) > 0 {
		if err := db.Model(&existing).Updates(updatedFields).Error; err != nil {
			fmt.Printf("Failed to update Http record: %v\n", err)
		} else {
			fmt.Printf("HTTP record updated for %s\n", http.SubDomain)
		}
	} else {
		fmt.Printf("No changes for %s, skipping update.\n", http.SubDomain)
	}
}

func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func compareMapFields(a, b MapField) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
