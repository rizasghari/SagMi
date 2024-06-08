package email

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/spf13/viper"
	"sagmi/internal/models"
	"time"
)

type MailgunService struct {
	Config *viper.Viper
	API    *mailgun.MailgunImpl
}

func NewMailgun(config *viper.Viper) *MailgunService {
	domain := config.GetString("mailgun.domain")
	apiKey := config.GetString("mailgun.api_key")
	regionIsEu := config.GetBool("mailgun.region_is_eu")
	mg := mailgun.NewMailgun(domain, apiKey)
	if regionIsEu {
		mg.SetAPIBase(mailgun.APIBaseEU)
	}
	return &MailgunService{
		Config: config,
		API:    mg,
	}
}

func (mgs *MailgunService) Send(alarmData models.AlarmData) error {
	sender := mgs.Config.GetString("mailgun.sender")
	receiver := mgs.Config.GetString("mailgun.receiver")
	subject := fmt.Sprintf("%s in %s is down!", alarmData.Service, alarmData.AppName)
	content := mgs.GetContent(alarmData)
	message := mgs.API.NewMessage(sender, subject, content, receiver)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, _, err := mgs.API.Send(ctx, message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (mgs *MailgunService) GetContent(alarmData models.AlarmData) string {
	return fmt.Sprintf(
		"Dear Team,\n\n"+
			"This is an alarm from Health Checker Service.\n"+
			"The %s service in %s app is down!\n"+
			"Log: %s\n\nSincerely,\nHealth Checker Service",
		alarmData.Service, alarmData.AppName,
		alarmData.Response,
	)
}
