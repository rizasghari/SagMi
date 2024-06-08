package app

import (
	"log"
	"sagmi/config"
	"sagmi/internal/cli"
	"sagmi/internal/controllers"
	"sagmi/internal/di"
	"sagmi/internal/repositories"
	"sagmi/internal/server"
	"sagmi/internal/services"
	"sync"
	"time"
)

type App struct {
	wg *sync.WaitGroup
}

func NewApp() *App {
	return &App{
		wg: &sync.WaitGroup{},
	}
}

func (a *App) LetsGo(restarted bool) {

	configuration := config.NewConfig()

	// Wire dependency injection used here
	db := di.InitDatabase()
	db.Start()

	logRepository := repositories.NewLogRepository(db.DB)
	logService := services.NewLoggerService(configuration.Viper, logRepository)

	endpointRepository := repositories.NewEndpointRepository(db.DB)
	endpointService := services.NewEndpointService(endpointRepository)

	healthCheckService := services.NewHealthCheckService(configuration.Viper, endpointService, logService)

	cliCheckCh := make(chan bool)
	_cli := cli.NewCLi(healthCheckService)
	go _cli.Run(cliCheckCh)

	alsoStarAutoCheckAndServer := <-cliCheckCh

	if alsoStarAutoCheckAndServer {
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			healthCheckService.StartAutomaticCheck()
		}()

		endpointController := controllers.NewEndpointController(endpointService, healthCheckService)

		httpServer := server.NewHttpServer(configuration.Viper, db.DB, endpointController)
		a.wg.Add(1)
		go func() {
			defer a.wg.Done()
			httpServer.Start()
		}()
	}

	// Defining the recover from panic mechanism
	restartAppAfter := time.Duration(configuration.Viper.GetInt("app.restart_in_after_panic")) * time.Second
	defer a.recoverFromPanic(restartAppAfter)

	a.wg.Wait()
}

func (a *App) recoverFromPanic(restartAppAfter time.Duration) {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic! Error: %v \nThe app will rstarted after %v...", r, restartAppAfter)
		time.Sleep(restartAppAfter)
		a.LetsGo(true)
	}
}
