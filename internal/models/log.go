package models

import (
	"database/sql"
	"time"
)

type Log struct {
	ID         uint         `json:"id,omitempty" gorm:"primaryKey"`
	EndpointId uint         `json:"endpoint_id,omitempty"`
	IsHealthy  bool         `json:"is_healthy"`
	Content    string       `json:"content"`
	DeletedAt  sql.NullTime `json:"deleted_at,omitempty"`
	CreatedAt  time.Time    `json:"created_at,omitempty"`
	UpdatedAt  time.Time    `json:"updated_at,omitempty"`
}
