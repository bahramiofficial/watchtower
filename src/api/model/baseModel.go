package model

import (
	"database/sql"
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
