//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"sagmi/config"
	"sagmi/internal/server"
)

func InitDatabase() *server.DatabaseServer {
	wire.Build(server.NewDatabaseServer, config.NewConfig)
	return &server.DatabaseServer{}
}
