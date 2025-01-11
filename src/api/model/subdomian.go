// Emails []string  `gorm:"type:text[];default:'{}'"`

package model

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bahramiofficial/watchtower/src/utilities"
	"gorm.io/gorm"
)

// /////////////////////////////////////////////////////////////
type Subdomain struct {
	BaseModel
	ProgramName string      `gorm:"type:text;not null;uniqueIndex:idx_program_subdomain"` // Named unique index
	SubDomain   string      `gorm:"type:text;not null;uniqueIndex:idx_program_subdomain"` // Same unique index name
	Scope       string      `gorm:"type:text; `
	Providers   StringArray `gorm:"type:text[]"`
}

// تعریف ایندکس ترکیبی
func (Subdomain) TableName() string {
	return "subdomains"
}

func (Subdomain) Indexes() []string {
	return []string{
		"UNIQUE (program_name, sub_domain)", // تعریف یکتایی برای ترکیب فیلدها
	}
}

type CreateUpdateSubDomainRequest struct {
	ProgramName string   `json:"programName" binding:"required"`
	SubDomain   string   `json:"subDomain" binding:"required"`
	Providers   []string `json:"providers"`
}

type SubDomainResponse struct {
	Id          int          `json:"id"`
	ProgramName string       `json:"programName"`
	SubDomain   string       `json:"subDomain"`
	Providers   []string     `json:"providers"`
	CreatedAt   time.Time    `json:"createdat"`
	UpdatedAt   sql.NullTime `json:"updatedat"`
}

// GetAllSubdomainWithScope fetches all subdomains associated with a given scope
func GetAllSubdomainWithScope(db *gorm.DB, scope string) ([]string, error) {
	// Initialize a variable to hold the results
	var subdomains []Subdomain

	// Fetch all subdomains where the scope matches the provided scope
	err := db.Where("scope = ?", scope).Find(&subdomains).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subdomains for scope '%s': %w", scope, err)
	}

	// Extract the subdomains from the result
	var subdomainNames []string
	for _, subdomain := range subdomains {
		subdomainNames = append(subdomainNames, subdomain.SubDomain)
	}

	// Return the list of subdomains
	return subdomainNames, nil
}

func FindSubdomainByProgramAndSubdomain(db *gorm.DB, programName, subDomain string) (*Subdomain, error) {

	// Initialize a variable to hold the result
	var subdomain Subdomain

	// Query the database
	err := db.Where("program_name = ? AND sub_domain = ?", programName, subDomain).First(&subdomain).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the record is not found, return nil and an error
			return nil, fmt.Errorf("subdomain '%s' for program '%s' not found", subDomain, programName)
		}
		// Return any other error encountered
		return nil, err
	}

	// If found, return the result and nil error
	return &subdomain, nil
}

// Always a single string for the provider
func UpsertSubdomain(db *gorm.DB, programName string, subDomain string, provider string) error {

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
	var subdomain Subdomain

	// Check if the subdomain exists
	err := db.Where("program_name = ? AND sub_domain = ?", programName, subDomain).First(&subdomain).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create a new subdomain if it doesn't exist
			newSubdomain := Subdomain{
				ProgramName: programName,
				SubDomain:   subDomain,
				Scope:       baseDomain,
				Providers:   StringArray{provider},
			}
			if err := db.Create(&newSubdomain).Error; err != nil {
				return fmt.Errorf("%v failed to insert subdomain: %w", utilities.GetFormattedTime(), err)
			}
			fmt.Printf("%v insert new subdomain: %v , with provider %v", utilities.GetFormattedTime(), subDomain, provider)
			return nil
		}
		return fmt.Errorf("%v failed to fetch subdomain: %w", utilities.GetFormattedTime(), err)
	}

	// Append the provider if it doesn't already exist
	for _, existingProvider := range subdomain.Providers {
		if existingProvider == provider {
			return nil // No update needed
		}
	}

	subdomain.Providers = append(subdomain.Providers, provider)

	// Save the updated subdomain
	if err := db.Save(&subdomain).Error; err != nil {
		return fmt.Errorf("%v failed to update subdomain: %w", utilities.GetFormattedTime(), err)
	}
	fmt.Printf("%v subdomain: %v  update  provider: %v", utilities.GetFormattedTime(), subDomain, provider)
	return nil
}
