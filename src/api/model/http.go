package model

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
