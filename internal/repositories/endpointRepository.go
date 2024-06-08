package repositories

import (
	"database/sql"
	"gorm.io/gorm"
	"sagmi/internal/models"
)

type EndpointRepository struct {
	db *gorm.DB
}

func NewEndpointRepository(dbHandler *gorm.DB) *EndpointRepository {
	return &EndpointRepository{db: dbHandler}
}

func (er *EndpointRepository) CreateNewEndpoint(endpoint *models.Endpoint) (*models.Endpoint, error) {
	result := er.db.Create(endpoint)
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return endpoint, nil
}

func (er *EndpointRepository) GetSingleEndpoint(id string) (*models.Endpoint, error) {
	var endpoint models.Endpoint
	result := er.db.Where("id = ?", id).First(&endpoint)
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &endpoint, nil
}

func (er *EndpointRepository) GetAllEndpoints() ([]*models.Endpoint, error) {
	var endpoints = make([]*models.Endpoint, 0)
	result := er.db.Find(&endpoints).Order("id desc")
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return endpoints, nil
}

func (er *EndpointRepository) GetAllEndpointsWithLatestLog() ([]*models.Result, error) {
	var results []*models.Result

	// Subquery to get the latest log for each endpoint
	subquery := er.db.Model(&models.Log{}).
		Select("id, endpoint_id, is_healthy, content, deleted_at, created_at").
		Where("id IN (?) AND deleted_at IS NULL",
			er.db.Table("logs").Select("MAX(id)").Group("endpoint_id"),
		)

	// Join with the subquery to get the latest log for each endpoint and include it in the results array
	rows, err := er.db.Table("endpoints").
		Select("endpoints.*, logs.*").
		Order("endpoints.id desc").
		Joins("left join (?) As logs on logs.endpoint_id = endpoints.id", subquery).
		Rows()

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var result models.Result
		if err := er.db.ScanRows(rows, &result); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return results, nil
}

func (er *EndpointRepository) DeleteEndpoint(id string) (*models.Endpoint, error) {
	return nil, nil
}
