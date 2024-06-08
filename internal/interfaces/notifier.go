package interfaces

import (
	"sagmi/internal/models"
)

type Notifier interface {
	Send(models.AlarmData) error
}
