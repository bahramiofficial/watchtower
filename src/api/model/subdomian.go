// Emails []string  `gorm:"type:text[];default:'{}'"`

package model

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bahramiofficial/watchtower/src/utilities"
	"gorm.io/gorm"
)

// StringArray custom type for PostgreSQL text[]
type StringArray []string

// Scan implements the sql.Scanner interface for StringArray
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = []string{}
		return nil
	}
	// Parse the PostgreSQL array literal
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to parse StringArray: %v", value)
	}
	// Convert the PostgreSQL array literal into a slice
	str = strings.Trim(str, "{}")
	if str == "" {
		*sa = []string{}
	} else {
		*sa = strings.Split(str, ",")
	}
	return nil
}

// Value implements the driver.Valuer interface for StringArray
func (sa StringArray) Value() (driver.Value, error) {
	// Convert the slice into a PostgreSQL-compatible array literal
	for i, v := range sa {
		sa[i] = fmt.Sprintf("\"%s\"", strings.ReplaceAll(v, "\"", "\\\""))
	}
	return "{" + strings.Join(sa, ",") + "}", nil
}

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
