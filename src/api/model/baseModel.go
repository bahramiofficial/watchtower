package model

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id        int          `gorm:"primaryKey"`
	CreatedAt time.Time    `gorm:"type:TIMESTAMP with time zone;default:CURRENT_TIMESTAMP; not null;"`
	UpdatedAt sql.NullTime `gorm:"type:TIMESTAMP with time zone;default:CURRENT_TIMESTAMP;  null;"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	m.CreatedAt = time.Now().UTC()
	return
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	//  id = tx.Statement.Context.Value("id")
	m.UpdatedAt = sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}
	return
}

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

// MapField type to handle JSONB fields in PostgreSQL
type MapField map[string]interface{}

// Scan implements the sql.Scanner interface for MapField
func (mf *MapField) Scan(value interface{}) error {
	if value == nil {
		*mf = make(MapField)
		return nil
	}
	str, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to parse MapField: %v", value)
	}
	return json.Unmarshal(str, mf)
}

// Value implements the driver.Valuer interface for MapField
func (mf MapField) Value() (driver.Value, error) {
	return json.Marshal(mf)
}
