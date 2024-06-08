package models

import (
	"database/sql"
	"time"
)

type Endpoint struct {
	ID            uint         `json:"id,omitempty" gorm:"primaryKey"`
	ParentAppName string       `json:"parent_app_name"`
	ServiceName   string       `json:"service_name"`
	Description   string       `json:"description"`
	Environment   string       `json:"environment"`
	HostName      string       `json:"host_name"`
	ServerIP      string       `json:"server_ip"`
	Path          string       `json:"path"`
	DeletedAt     sql.NullTime `json:"deleted_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at,omitempty"`
	UpdatedAt     time.Time    `json:"updated_at,omitempty"`
}
