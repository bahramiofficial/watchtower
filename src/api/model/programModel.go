// Emails []string  `gorm:"type:text[];default:'{}'"`

package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

type ProgramModel struct {
	BaseModel
	ProgramName string          `gorm:"type:text;not null;uniqueIndex"` // Unique and indexed
	Config      json.RawMessage `gorm:"type:jsonb;default:null"`        // Dictionary field, nullable
	Scopes      []string        `gorm:"type:text[];default:null"`       // Array field, nullable
	Otoscopes   []string        `gorm:"type:text[];default:null"`       // Array field, nullable
}

type CreateUpdateProgramModelRequest struct {
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

type ProgramModelResponse struct {
	Id          int             `json:"id"`
	ProgramName string          `json:"programName"`
	Config      json.RawMessage `json:"config"`
	Otoscopes    []string        `json:"otoscopes"`
	CreatedAt   time.Time       `json:"createdat"`
	UpdatedAt   sql.NullTime    `json:"updatedat"`
}

// /////////////////////////////////////////////////////////////
type SubDomainModel struct {
	BaseModel
	ProgramName string   `gorm:"type:text;not null;uniqueIndex"` // Unique and indexed
	SubDomain   string   `gorm:"type:text;not null;uniqueIndex"` // Unique and indexed
	Providers   []string `gorm:"type:text[]"`                    // Array field, no default value

}

type CreateUpdateSubDomainModelRequest struct {
	ProgramName string   `json:"programName" binding:"required"`
	SubDomain   string   `json:"subDomain" binding:"required"`
	Providers   []string `json:"providers"`
}

type SubDomainModelResponse struct {
	Id          int          `json:"id"`
	ProgramName string       `json:"programName"`
	SubDomain   string       `json:"subDomain"`
	Providers   []string     `json:"providers"`
	CreatedAt   time.Time    `json:"createdat"`
	UpdatedAt   sql.NullTime `json:"updatedat"`
}
