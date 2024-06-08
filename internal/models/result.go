package models

type Result struct {
	Endpoint
	Log Log `gorm:"embedded"`
}
