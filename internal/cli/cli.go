package cli

import (
	"flag"
	"fmt"
	"sagmi/internal/models"
	"sagmi/internal/services"
)

type Cli struct {
	healthCheckService *services.HealthCheckService
}

func NewCLi(healthCheckService *services.HealthCheckService) *Cli {
	return &Cli{
		healthCheckService: healthCheckService,
	}
}

func (c *Cli) Run(alsoStarAutoCheckAndServerCh chan bool) {

	var hostName string
	var serverIP string
	var path string
	var alsoStarAutoCheckAndServer bool

	flag.StringVar(&hostName, "host", "", "Endpoint host name")
	flag.StringVar(&serverIP, "ip", "", "Endpoint server IP address")
	flag.StringVar(&path, "path", "", "Endpoint path")
	flag.BoolVar(&alsoStarAutoCheckAndServer, "auto", false, "Also star auto check and http server")

	flag.Parse()

	if hostName != "" || serverIP != "" || path != "" {
		if hostName == "" || serverIP == "" || path == "" {
			fmt.Println("Please provide host name, server IP address and path. Example: -host=example.com -ip=1.1.1.1 -path=/")
			alsoStarAutoCheckAndServerCh <- false
		} else {
			endpoint := &models.Endpoint{
				HostName: hostName,
				ServerIP: serverIP,
				Path:     path,
			}
			err := c.healthCheckService.CliHealthCheck(endpoint)
			if err != nil {
				fmt.Printf("CliHealthCheck for %v%v - ERROR: %v\n", endpoint.HostName, endpoint.Path, err)
				alsoStarAutoCheckAndServerCh <- alsoStarAutoCheckAndServer
			} else {
				fmt.Printf("CliHealthCheck for %v%v - OK\n", endpoint.HostName, endpoint.Path)
				alsoStarAutoCheckAndServerCh <- alsoStarAutoCheckAndServer
			}
		}
	} else {
		alsoStarAutoCheckAndServerCh <- true
	}
	return
}
