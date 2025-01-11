package model

// LiveSubdomains represents the live subdomains model
type LiveSubdomains struct {
	BaseModel
	ProgramName string      `gorm:"type:text;not null;uniqueIndex:idx_program_subdomain"` // Composite unique index
	SubDomain   string      `gorm:"type:text;not null;uniqueIndex:idx_program_subdomain"` // Same unique index name
	Scope       string      `gorm:"type:text;not null"`
	Cdn         string      `gorm:"type:text;"`
	IPs         StringArray `gorm:"type:text[]"`
	Tag         string      `gorm:"type:text"`
}
 