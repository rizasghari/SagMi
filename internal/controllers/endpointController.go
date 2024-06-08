package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sagmi/internal/models"
	"sagmi/internal/services"
	"time"
)

type EndpointController struct {
	endpointService    *services.EndpointService
	healthCheckService *services.HealthCheckService
}

func NewEndpointController(endpointService *services.EndpointService, healthCheckService *services.HealthCheckService) *EndpointController {
	return &EndpointController{
		endpointService:    endpointService,
		healthCheckService: healthCheckService,
	}
}

func (ec *EndpointController) CreateNewEndpoint(ctx *gin.Context) {
	var endpoint models.Endpoint
	endpoint.ParentAppName = ctx.PostForm("parent_app_name")
	endpoint.ServiceName = ctx.PostForm("service_name")
	endpoint.Description = ctx.PostForm("description")
	endpoint.Environment = ctx.PostForm("environment")
	endpoint.HostName = ctx.PostForm("host_name")
	endpoint.ServerIP = ctx.PostForm("server_ip")
	endpoint.Path = ctx.PostForm("path")
	ec.endpointService.CreateNewEndpoint(&endpoint)
}

func (ec *EndpointController) GetSingleEndpoint(ctx *gin.Context) {
	id := ctx.Param("id")
	endpoint, err := ec.endpointService.GetSingleEndpoint(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, endpoint)
}

func (ec *EndpointController) GetAllEndpointsWithLatestLog(ctx *gin.Context) {
	ec.healthCheckService.ManualHealthCheck()
	endpoints, err := ec.endpointService.GetAllEndpointsWithLatestLog()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	time.Sleep(1 * time.Second)
	ctx.HTML(http.StatusOK, "endpoints.html", endpoints)
}

func (ec *EndpointController) GetIndexPage(ctx *gin.Context) {
	results, err := ec.endpointService.GetAllEndpointsWithLatestLog()
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	content := models.Content{Results: results, PageTitle: "API Health Checker Service"}
	ctx.HTML(http.StatusOK, "index.html", content)
}

func (ec *EndpointController) DeleteEndpoint(ctx *gin.Context) {
	id := ctx.Param("id")
	ec.endpointService.DeleteEndpoint(id)
}
