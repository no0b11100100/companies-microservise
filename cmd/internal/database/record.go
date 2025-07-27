package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CompanyInfo represents a company object
type CompanyInfo struct {
	ID             *uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name           *string    `json:"name" gorm:"size:15;not null;uniqueIndex"`
	Description    *string    `json:"description,omitempty" gorm:"size:3000"`
	EmployeesCount *int       `json:"employeesCount" gorm:"not null"`
	IsRegistered   *bool      `json:"isRegistered" gorm:"not null"`
	Type           *int       `json:"type" gorm:"not null"`
}

func (r *CompanyInfo) BeforeCreate(tx *gorm.DB) error {
	if r.ID == nil {
		id := uuid.New()
		r.ID = &id
	}
	return nil
}
