package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CompanyInfo represents a company object
type CompanyInfo struct {
	CompanyID    *uuid.UUID `json:"companyid,omitempty" gorm:"type:char(36);not null;uniqueIndex"`
	Name         *string    `json:"name"`
	Description  *string    `json:"description,omitempty"`
	Employees    *int       `json:"employees"`
	IsRegistered *bool      `json:"isRegistered"`
	Type         *int       `json:"type"`
}

type Record struct {
	ID int `gorm:"primaryKey"`
	CompanyInfo
}

func (r *Record) BeforeCreate(tx *gorm.DB) error {
	if r.CompanyID == nil {
		id := uuid.New()
		r.CompanyID = &id
	}
	return nil
}
