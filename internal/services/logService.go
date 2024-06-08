package services

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"sagmi/internal/models"
	"sagmi/internal/repositories"
	"strings"
	"time"
)

type LoggerService struct {
	Config        *viper.Viper
	logRepository *repositories.LogRepository
}

func NewLoggerService(config *viper.Viper, logRepository *repositories.LogRepository) *LoggerService {
	return &LoggerService{
		Config:        config,
		logRepository: logRepository,
	}
}

func (ls *LoggerService) SaveLogInFile(message string) error {
	logFilePath := ls.Config.GetString("logger.log_file_path")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	message = currentTime + " " + strings.Trim(message, "\t") + "\n"
	if _, err := logFile.Write([]byte(message)); err != nil {
		log.Println(err)
	}
	if err := logFile.Close(); err != nil {
		log.Println(err)
	}
	return nil
}

func (ls *LoggerService) SaveLogInDatabase(log *models.Log) error {
	_, err := ls.logRepository.Create(log)
	if err != nil {
		return err
	}
	return nil
}

func (ls *LoggerService) PurgeLogs() {
	// Todo: Implement auto periodic log purge to remove old logs
}
