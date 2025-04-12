package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string    `gorm:"unique;not null"`
	Password      string    `gorm:"unique;not null"`
	Server        string    `gorm:"unique;not null"`
	LastOnline    time.Time `gorm:"not null"`
	StoreMessages bool      `gorm:"not null"`
	StoreTimeline bool      `gorm:"not null"`
}
