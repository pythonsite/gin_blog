package models

import "time"

type BaseModel struct {
	ID uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// query result
type QrArchive struct {
	ArchiveDate time.Time //month
	Total       int       //total
	Year        int       // year
	Month       int       // month
}