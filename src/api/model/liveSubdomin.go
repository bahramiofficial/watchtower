package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bahramiofficial/watchtower/src/utilities"
	"gorm.io/gorm"
)

// LiveSubdomains represents the live subdomains model
type LiveSubdomain struct {
	BaseModel
	ProgramName string      `gorm:"type:text;not null;uniqueIndex:idx_program_subdomain"` // Composite unique index
	SubDomain   string      `gorm:"type:text;not null;uniqueIndex:idx_program_subdomain"` // Same unique index name
	Scope       string      `gorm:"type:text;not null"`
	Cdn         string      `gorm:"type:text;"`
	IPs         StringArray `gorm:"type:text[]"`
	Tag         string      `gorm:"type:text"`
}

// GetAllLiveSubdomainWithScope fetches all subdomains associated with a given scope
func GetAllLiveSubdomainWithScope(db *gorm.DB, scope string) ([]LiveSubdomain, error) {
	// Initialize a variable to hold the results
	var liveSubdomains []LiveSubdomain

	// Fetch all subdomains where the scope matches the provided scope
	err := db.Where("scope = ?", scope).Find(&liveSubdomains).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subdomains for scope '%s': %w", scope, err)
	}

	// Return the list of subdomains
	return liveSubdomains, nil
}

// UpsertLiveSubdomain creates or updates a live subdomain entry in the database
func UpsertLiveSubdomain(db *gorm.DB, programName string, subDomain string, scope string, ips []string, tag string) error {

	// Check if the subdomain matches the domain and is not a wildcard
	// Ensure subdomain is valid (no wildcard and no top-level domain)
	if strings.Contains(subDomain, "*") {
		return fmt.Errorf("subdomain '%s' contains a wildcard (*), which is not allowed", subDomain)
	}

	// Check if the subdomain matches the domain and is not a wildcard
	if strings.Count(subDomain, ".") <= 1 {
		return fmt.Errorf("subdomain '%s' is invalid. It must contain at least one subdomain (e.g., sub.x.com)", subDomain)
	}

	// Fetch the program using the program name
	var program Program
	if err := db.Where("program_name = ?", programName).First(&program).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%v program '%s' not found", utilities.GetFormattedTime(), programName)
		}
		return fmt.Errorf("failed to fetch program: %w", err)
	}

	// Extract the base domain from the subdomain
	baseDomain := utilities.ExtractBaseDomain(subDomain)

	// Check if the base domain is in the program's Scopes or Otoscopes
	inScope := false
	for _, scope := range program.Scopes {
		if baseDomain == scope {
			inScope = true
			break
		}
	}
	for _, oto := range program.Otoscopes {
		if baseDomain == oto {
			inScope = false
			break
		}
	}

	// Check if the base domain is in the program's Scopes or Otoscopes

	for _, scope := range program.Scopes {
		if subDomain == scope {
			inScope = true
			break
		}
	}
	for _, oto := range program.Otoscopes {
		if subDomain == oto {
			inScope = false
			break
		}
	}

	// If the base domain is not in Scopes or Otoscopes, print a warning and return
	if !inScope {
		fmt.Printf("%v Subdomain '%s' not in scope for program '%s'\n", utilities.GetFormattedTime(), subDomain, programName)
		return nil
	}

	// Initialize a variable to hold the result
	var liveSubdomain LiveSubdomain

	// Check if the live subdomain exists
	err := db.Where("program_name = ? AND sub_domain = ?", programName, subDomain).First(&liveSubdomain).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a new live subdomain if it doesn't exist
			newLiveSubdomain := LiveSubdomain{
				ProgramName: programName,
				SubDomain:   subDomain,
				Scope:       scope,
				Cdn:         "",
				IPs:         ips,
				Tag:         tag,
			}
			if err := db.Create(&newLiveSubdomain).Error; err != nil {
				return fmt.Errorf("%v failed to insert live subdomain: %w", utilities.GetFormattedTime(), err)
			}
			fmt.Printf("%v insert new live subdomain: %v\n", utilities.GetFormattedTime(), subDomain)
			return nil
		}
		return fmt.Errorf("%v failed to fetch live subdomain: %w", utilities.GetFormattedTime(), err)
	}

	// Update the IPs if they are not already included
	ipMap := make(map[string]bool)
	for _, existingIP := range liveSubdomain.IPs {
		ipMap[existingIP] = true
	}
	for _, ip := range ips {
		if !ipMap[ip] {
			liveSubdomain.IPs = append(liveSubdomain.IPs, ip)
		}
	}

	// Update the tag if provided
	if tag != "" {
		liveSubdomain.Tag = tag
	}

	// Save the updated live subdomain
	if err := db.Save(&liveSubdomain).Error; err != nil {
		return fmt.Errorf("%v failed to update live subdomain: %w", utilities.GetFormattedTime(), err)
	}
	fmt.Printf("%v updated live subdomain: %v\n", utilities.GetFormattedTime(), subDomain)
	return nil
}
