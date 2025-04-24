package model

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID           int            `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	IsProject    bool           `json:"is_project"`
	ClientPrefix string         `json:"client_prefix"`
	ClientLogo   string         `json:"client_logo"`
	Address      string         `json:"address"`
	PhoneNumber  string         `json:"phone_number"`
	City         string         `json:"city"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
}
