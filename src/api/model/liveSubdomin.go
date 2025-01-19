package model

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bahramiofficial/watchtower/src/utilities"
	"gorm.io/gorm"
)

// LiveSubdomains represents the live subdomains model
type LiveSubdomain struct {
	BaseModel
	ProgramName string      `gorm:"type:text;not null;uniqueIndex:idx_program_subdomain"` // Composite unique index
	Subdomain   string      `gorm:"type:text;not null;uniqueIndex:idx_program_subdomain"` // Same unique index name
	Scope       string      `gorm:"type:text;not null"`
	IPs         StringArray `gorm:"type:text[]"`
	Tag         string      `gorm:"type:text"`
}

// GetLiveSubdomain fetches a LiveSubdomain record by subdomain.
func GetLiveSubdomain(db *gorm.DB, subdomain string) (*LiveSubdomain, error) {
	var liveSubdomain LiveSubdomain
	err := db.Where("subdomain = ?", subdomain).First(&liveSubdomain).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err // Return nil without an error if no record is found
		}
		return nil, fmt.Errorf("failed to fetch live subdomain: %w", err)
	}
	return &liveSubdomain, nil
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

func GetAllLiveSubdomainWithScopeName(db *gorm.DB, scope string) ([]string, error) {
	// Initialize a variable to hold the results
	var liveSubdomain []LiveSubdomain

	// Fetch all subdomains where the scope matches the provided scope
	err := db.Where("scope = ?", scope).Find(&liveSubdomain).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subdomains for scope '%s': %w", scope, err)
	}

	// Extract the subdomains from the result
	var livesubdomainNames []string
	for _, subdomain := range liveSubdomain {
		livesubdomainNames = append(livesubdomainNames, subdomain.Subdomain)
	}

	// Return the list of subdomains
	return livesubdomainNames, nil
}

// GetAllLiveSubdomainWithScope fetches all subdomains associated with a given scope
func GetAllLiveSubdomainWithScopeAndDomain(db *gorm.DB, scope string, domain string) ([]LiveSubdomain, error) {
	// Initialize a variable to hold the results
	var liveSubdomains []LiveSubdomain

	// Fetch all subdomains where the scope matches the provided scope
	err := db.Where("scope = ? AND ", scope).Find(&liveSubdomains).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subdomains for scope '%s': %w", scope, err)
	}

	// Return the list of subdomains
	return liveSubdomains, nil
}

// UpsertLives inserts or updates a live subdomain entry in the database.
func UpsertLiveSubdomain(db *gorm.DB, domain, subdomain string, ips []string, tag string) error {
	subdomain = strings.ToLower(subdomain) // Ensure subdomain is in lowercase

	// Get the program associated with the domain
	program, err := GetProgramByScope(db, domain)
	if err != nil {
		return fmt.Errorf("failed to fetch program for domain %s: %w", domain, err)
	}

	// Fetch an existing live subdomain if it exists
	existing, _ := GetLiveSubdomain(db, subdomain)

	// If the subdomain exists, update its IPs if needed
	if existing != nil {
		print(existing.ProgramName, "thi isssssssss")
		sort.Strings(existing.IPs) // Sort existing IPs
		sort.Strings(ips)          // Sort the new IPs
		if !utilities.AreSlicesEqual(existing.IPs, ips) {
			existing.IPs = ips
			log.Printf("[%s] Updated live subdomain: %s", utilities.GetFormattedTime(), subdomain)
		}

		err = db.Save(&existing).Error
		if err != nil {
			return fmt.Errorf("failed to update subdomain: %w", err)
		}
	} else {
		print("sskhadkshdlkjahdsthi isssssssss")
		// Insert a new live subdomain if none exists
		newLiveSubdomain := &LiveSubdomain{
			ProgramName: program.ProgramName,
			Subdomain:   subdomain,
			Scope:       domain,
			IPs:         ips,
			Tag:         tag,
		}
		err := db.Save(&newLiveSubdomain).Error
		if err != nil {
			return fmt.Errorf("failed to save new subdomain: %w", err)
		}
		utilities.SendDiscordMessage(fmt.Sprintf("```'%s' (fresh live) has been added to '%s' program```", subdomain, program.ProgramName))

	}

	return nil
}
