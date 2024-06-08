package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"sagmi/internal/controllers"
)

type HttpServer struct {
	Config             *viper.Viper
	Router             *gin.Engine
	db                 *gorm.DB
	endpointController *controllers.EndpointController
}

func NewHttpServer(config *viper.Viper, dbHandler *gorm.DB, endpointController *controllers.EndpointController) *HttpServer {
	return &HttpServer{
		Config:             config,
		db:                 dbHandler,
		endpointController: endpointController,
	}
}

func (hs *HttpServer) Start() {
	hs.Router = gin.Default()

	hs.Router.Static("/web/static/", "./web/static")
	hs.Router.LoadHTMLGlob("web/templates/*")

	hs.Router.GET("/", hs.endpointController.GetIndexPage)
	hs.Router.POST("/endpoint", hs.endpointController.CreateNewEndpoint)
	hs.Router.GET("/endpoint", hs.endpointController.GetAllEndpointsWithLatestLog)
	hs.Router.GET("/endpoint/:id", hs.endpointController.GetSingleEndpoint)
	hs.Router.DELETE("/endpoint/:id", hs.endpointController.DeleteEndpoint)

	hs.Run()
}

func (hs *HttpServer) Run() {
	port := fmt.Sprintf(":%s", hs.Config.GetString("app.port"))
	err := hs.Router.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
