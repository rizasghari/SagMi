package messenger

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"net/http"
	"sagmi/internal/models"
	"time"
)

type SlackService struct {
	Config *viper.Viper
}

func NewSlack(config *viper.Viper) *SlackService {
	return &SlackService{Config: config}
}

func (s *SlackService) getMessageBody(params models.AlarmData) *models.MessageBody {
	return &models.MessageBody{
		Text:     "<!channel> " + params.Service + " in " + params.AppName + " is down!",
		Markdown: true,
		UserName: "Health Check",
		Attachments: []models.Attachment{
			{
				Fields: []models.Field{
					{
						Short: true,
						Title: "Health Check Service",
						Value: params.HealthCheckService,
					},
					{
						Short: true,
						Title: "Environment",
						Value: params.Environment,
					},
					{
						Short: true,
						Title: "App Name",
						Value: params.AppName,
					},
					{
						Short: true,
						Title: "Service",
						Value: params.Service,
					},
					{
						Short: true,
						Title: "Date & Time",
						Value: time.Now().Format("2006-01-02 15:04:05"),
					},
					{
						Short: true,
						Title: "Response",
						Value: params.Response,
					},
				},
				Actions: []models.Action{
					{
						Type: "button",
						Text: "Health Check",
						Url:  params.ServiceURL,
					},
					{
						Type: "button",
						Text: "Workloads",
						Url:  s.Config.GetString("google_cloud.production_project_url"),
					},
				},
			},
		},
	}
}

func (s *SlackService) Send(params models.AlarmData) error {
	body := s.getMessageBody(params)
	webHookUrl := s.Config.GetString("slack.webhook_url")

	if webHookUrl == "" {
		return errors.New("webhook url is empty")
	}

	client := &http.Client{}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webHookUrl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
