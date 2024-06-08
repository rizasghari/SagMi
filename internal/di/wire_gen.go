// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"sagmi/config"
	"sagmi/internal/server"
)

// Injectors from wire.go:

func InitDatabase() *server.DatabaseServer {
	configConfig := config.NewConfig()
	databaseServer := server.NewDatabaseServer(configConfig)
	return databaseServer
}