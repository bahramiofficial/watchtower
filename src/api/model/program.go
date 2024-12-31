// Emails []string  `gorm:"type:text[];default:'{}'"`

package model

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// Import pq package for PostgreSQL support
type Program struct {
	BaseModel
	ProgramName string          `gorm:"type:text;not null;uniqueIndex"` // Unique and indexed
	Config      json.RawMessage `gorm:"type:jsonb;default:null"`        // Dictionary field, nullable
	Scopes      pq.StringArray  `gorm:"type:text[]"`                    // Correctly handle PostgreSQL text[]
	Otoscopes   pq.StringArray  `gorm:"type:text[]"`                    // Correctly handle PostgreSQL text[]
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

// /////////////////////////////////////////////////////////////
type Subdomain struct {
	BaseModel
	ProgramName string   `gorm:"type:text;not null;uniqueIndex"` // Unique and indexed
	SubDomain   string   `gorm:"type:text;not null;uniqueIndex"` // Unique and indexed
	Providers   []string `gorm:"type:text[]"`                    // Array field, no default value

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
