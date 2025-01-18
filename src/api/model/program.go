// Emails []string  `gorm:"type:text[];default:'{}'"`

package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Import pq package for PostgreSQL support
type Program struct {
	BaseModel
	ProgramName string          `gorm:"type:text;not null;uniqueIndex"` // Unique and indexed
	Config      json.RawMessage `gorm:"type:jsonb;default:null"`        // Dictionary field, nullable
	Scopes      pq.StringArray  `gorm:"type:text[]"`                    // Correctly handle PostgreSQL text[]
	Otoscopes   pq.StringArray  `gorm:"type:text[]"`                    // Correctly handle PostgreSQL text[]

}

// getAllPrograms retrieves all programs from the database
func GetAllPrograms(db *gorm.DB) ([]Program, error) {
	var programs []Program
	if err := db.Find(&programs).Error; err != nil {
		return nil, err
	}
	return programs, nil
}

// GetProgramByProgramName fetches a Program by its ProgramName
func GetProgramByProgramName(db *gorm.DB, programName string) (Program, error) {
	var program Program

	// Fetch the program where ProgramName matches the provided value
	err := db.Where("program_name = ?", programName).First(&program).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return program, fmt.Errorf("no program found with program_name: %s", programName)
		}
		return program, fmt.Errorf("failed to fetch program by program_name: %w", err)
	}

	return program, nil
}

// FindDomainWithProgramName queries a Program by its ProgramName
func FindDomainWithProgramName(db *gorm.DB, programName string) (*Program, error) {
	// Initialize a variable to hold the result
	var program Program

	// Query the database
	err := db.Where("program_name = ?", programName).First(&program).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the record is not found, return nil and an error
			return nil, fmt.Errorf("program with name '%s' not found", programName)
		}
		// Return any other error encountered
		return nil, err
	}

	// If found, return the result and nil error
	return &program, nil
}

// insertOrUpdateProgram inserts a new record or updates an existing one based on the ProgramName.
func InsertOrUpdateProgram(db *gorm.DB, program *Program) error {
	// Check if the program already exists
	var existing Program
	err := db.Where("program_name = ?", program.ProgramName).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Record not found, create a new one
			return db.Create(program).Error
		}
		// Other errors
		return err
	}

	// Record found, update it
	existing.Config = program.Config
	existing.Scopes = program.Scopes
	existing.Otoscopes = program.Otoscopes

	return db.Save(&existing).Error
}

type CreateUpdateProgramRequest struct {
	ProgramName string          `json:"name" binding:"required,min:3,max:100"`
	Config      json.RawMessage `json:"config"`
	Scopes      []string        `json:"scopes"`
	Otoscopes   []string        `json:"otoscopes"`
}

// type UpdateProgramModelRequest struct {
// 	Name      string          `json:"name" binding:"required,min:3,max:100"`
// 	Config    json.RawMessage `json:"config" binding`
// 	Scopes    []string        `json:"scopes" binding`
// 	Otoscopes []string        `json:"ooscopes" binding`
// }

type ProgramResponse struct {
	Id          int             `json:"id"`
	ProgramName string          `json:"programName"`
	Config      json.RawMessage `json:"config"`
	Otoscopes   []string        `json:"otoscopes"`
	CreatedAt   time.Time       `json:"createdat"`
	UpdatedAt   sql.NullTime    `json:"updatedat"`
}

func AddNewProgramIfNotExist(db *gorm.DB, programName string, config json.RawMessage, scopes, otoscopes []string) (*int, error) {
	// Check if the program already exists
	var program Program
	err := db.Where("program_name = ?", programName).First(&program).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If the program doesn't exist, create a new one
			newProgram := Program{
				ProgramName: programName,
				Config:      config,
				Scopes:      pq.StringArray(scopes),
				Otoscopes:   pq.StringArray(otoscopes),
			}

			// Insert the new program into the database
			if err := db.Create(&newProgram).Error; err != nil {
				// Return an error if insertion fails
				return nil, fmt.Errorf("failed to insert program: %w", err)
			}

			// Return the new program ID
			return &newProgram.Id, nil
		}
		// Return other errors encountered during the query
		return nil, fmt.Errorf("failed to check if program exists: %w", err)
	}

	// If the program already exists, return nil and an error
	return nil, fmt.Errorf("program with name '%s' already exists", programName)
}

// use abuse db
// GetProgramByScope retrieves the first program that matches the given scope.
func GetProgramByScope(db *gorm.DB, domain string) (*Program, error) {
	var program Program
	err := db.Where("ARRAY_TO_STRING(scopes, ',') ILIKE ?", "%"+domain+"%").First(&program).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no program found for scope '%s'", domain)
		}
		return nil, err
	}
	return &program, nil
}
