package repositories

import (
	"gorm.io/gorm"
	"sagmi/internal/models"
	"time"
)

type LogRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{db: db}
}

func (lr *LogRepository) Create(log *models.Log) (*models.Log, error) {
	result := lr.db.Create(log)
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return log, nil
}

func (lr *LogRepository) GetAll() ([]*models.Log, error) {
	var logs []*models.Log
	result := lr.db.Find(&logs)
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return logs, nil
}

func (lr *LogRepository) GetAllByEndpointId(endpointId uint) ([]*models.Log, error) {
	var logs []*models.Log
	result := lr.db.Where("endpoint_id = ?", endpointId).Find(&logs)
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return logs, nil
}

func (lr *LogRepository) GetEndpointLastLog(endpointId uint) (*models.Log, error) {
	var log *models.Log
	result := lr.db.Where("endpoint_id = ?", endpointId).Order("created_at desc").First(&log)
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return log, nil
}

func (lr *LogRepository) SoftDelete(log *models.Log) (*models.Log, error) {
	result := lr.db.Model(&log).Update("deleted_at", time.Now())
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return log, nil
}

func (lr *LogRepository) HardDelete(log *models.Log) error {
	result := lr.db.Delete(&log)
	if err := result.Error; err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
