package services

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"sagmi/internal/interfaces"
	"sagmi/internal/models"
	"sagmi/internal/services/notifiers/email"
	"sagmi/internal/services/notifiers/messenger"
	"sync"
	"time"
)

var (
	Unsuccessful = errors.New("health check request failed")
	NotFound     = errors.New("404 - The service URL is not found")
)

type HealthCheckService struct {
	Config          *viper.Viper
	Slack           interfaces.Notifier
	Email           interfaces.Notifier
	logService      *LoggerService
	endpointService *EndpointService
	httpClient      *http.Client
}

func NewHealthCheckService(config *viper.Viper, endpointService *EndpointService, logService *LoggerService) *HealthCheckService {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: transport}
	httpClient.Timeout = time.Duration(config.GetInt("app.request_timeout")) * time.Second

	slack := messenger.NewSlack(config)
	mailgun := email.NewMailgun(config)

	return &HealthCheckService{
		Config:          config,
		Slack:           slack,
		Email:           mailgun,
		logService:      logService,
		endpointService: endpointService,
		httpClient:      httpClient,
	}
}

func (hcs *HealthCheckService) LoadEndpointsFromJson() ([]models.Endpoint, error) {
	endpoints := make([]models.Endpoint, 0)
	jsonFile, err := os.Open("./static/json/apps.json")
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(jsonFile).Decode(&endpoints)
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (hcs *HealthCheckService) LoadEndpointsFromDB() ([]*models.Endpoint, error) {
	return hcs.endpointService.GetAllEndpoints()
}

func (hcs *HealthCheckService) CliHealthCheck(endpoint *models.Endpoint) error {
	healthCheckURl := fmt.Sprintf("https://%v%v", endpoint.ServerIP, endpoint.Path)
	return hcs.Request(endpoint, healthCheckURl)
}

func (hcs *HealthCheckService) ManualHealthCheck() {
	wg := &sync.WaitGroup{}
	endpoints, err := hcs.LoadEndpointsFromDB()
	if err != nil {
		log.Println(err)
	}
	log.Println("🔄Running manual health check")
	for _, endpoint := range endpoints {
		wg.Add(1)
		go hcs.Check(endpoint, true, wg)
	}
	wg.Wait()
}

func (hcs *HealthCheckService) StartAutomaticCheck() {
	for {
		wg := &sync.WaitGroup{}

		endpoints, err := hcs.LoadEndpointsFromDB()
		if err != nil {
			log.Println(err)
		}

		log.Println("🔄Running health check round")
		for _, endpoint := range endpoints {
			wg.Add(1)
			go hcs.Check(endpoint, false, wg)
		}

		runEvery := time.Duration(hcs.Config.GetInt("app.run_every_time")) * time.Second
		// To ensure app performance is not affected by the time it takes to run the health check
		if runEvery < 30*time.Second {
			runEvery = 30 * time.Second
		}
		time.Sleep(runEvery)

		wg.Wait()
	}
}

func (hcs *HealthCheckService) Check(endpoint *models.Endpoint, isManual bool, wg *sync.WaitGroup) {
	healthCheckURl := fmt.Sprintf("https://%v%v", endpoint.ServerIP, endpoint.Path)
	healthCheckError := hcs.Request(endpoint, healthCheckURl)
	var message string
	newLog := &models.Log{
		EndpointId: endpoint.ID,
		IsHealthy:  healthCheckError == nil,
	}
	if healthCheckError != nil {
		newLog.Content = healthCheckError.Error()
		message = fmt.Sprintf("\t❌ %v in %v is down! Error: %v", endpoint.ServiceName, endpoint.ParentAppName, healthCheckError)

		// Send health check alarm to the Slack channel if is active and not manual check
		sendAlarmToSlack := hcs.Config.GetBool("slack.send_alarm")
		if sendAlarmToSlack && !isManual {
			params := models.AlarmData{
				Environment:        endpoint.Environment,
				AppName:            endpoint.ParentAppName,
				Service:            endpoint.ServiceName,
				ServiceURL:         healthCheckURl,
				Response:           healthCheckError.Error(),
				HealthCheckService: hcs.Config.GetString("app.name"),
			}
			err := hcs.Slack.Send(params)
			if err != nil {
				return
			}
		}

		// Send email if is active and not manual check
		sendEmail := hcs.Config.GetBool("mailgun.send_email")
		if sendEmail && !isManual {
			alarmData := models.AlarmData{
				Environment:        endpoint.Environment,
				AppName:            endpoint.ParentAppName,
				Service:            endpoint.ServiceName,
				ServiceURL:         healthCheckURl,
				Response:           healthCheckError.Error(),
				HealthCheckService: hcs.Config.GetString("app.name"),
			}
			err := hcs.Email.Send(alarmData)
			if err != nil {
				return
			}
		}
	} else {
		message = fmt.Sprintf("\t✅ %v in %v is healthy", endpoint.ServiceName, endpoint.ParentAppName)
	}

	// Save new log in database if database log is active
	saveLogInDb := hcs.Config.GetBool("logger.save_in_database")
	if saveLogInDb /*&& !newLog.IsHealthy*/ {
		err := hcs.logService.SaveLogInDatabase(newLog)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// Print logs in console if is not manual and is console log is active
	logInConsole := hcs.Config.GetBool("logger.print_in_console")
	if logInConsole && !isManual {
		log.Println(message)
	}

	// Save logs in file if is not manual and file log is active
	saveLogInFile := hcs.Config.GetBool("logger.save_in_file")
	if saveLogInFile && !newLog.IsHealthy && !isManual {
		err := hcs.logService.SaveLogInFile(message)
		if err != nil {
			log.Println(err.Error())
		}
	}

	wg.Done()
}

func (hcs *HealthCheckService) Request(endpoint *models.Endpoint, healthCheckURl string) error {
	req, err := http.NewRequest("GET", healthCheckURl, nil)
	if err != nil {
		return err
	}
	req.Host = endpoint.HostName
	resp, err := hcs.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode == 404 {
		return NotFound
	}

	if resp.StatusCode != 200 {
		return Unsuccessful
	}

	return nil
}
