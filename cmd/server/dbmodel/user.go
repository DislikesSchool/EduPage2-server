package dbmodel

import (
	"time"
)

type User struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Username      string    `gorm:"unique;not null"`
	Password      string    `gorm:"unique;not null"`
	Server        string    `gorm:"unique;not null"`
	LastOnline    time.Time `gorm:"not null"`
	StoreMessages bool      `gorm:"not null"`
	StoreTimeline bool      `gorm:"not null"`
}
