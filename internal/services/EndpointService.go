package services

import (
	"sagmi/internal/models"
	"sagmi/internal/repositories"
)

type EndpointService struct {
	EndpointRepository *repositories.EndpointRepository
}

func NewEndpointService(endpointRepository *repositories.EndpointRepository) *EndpointService {
	return &EndpointService{EndpointRepository: endpointRepository}
}

func (e *EndpointService) CreateNewEndpoint(endpoint *models.Endpoint) {
	_, err := e.EndpointRepository.CreateNewEndpoint(endpoint)
	if err != nil {
		return
	}
}

func (e *EndpointService) GetSingleEndpoint(id string) (*models.Endpoint, error) {
	return e.EndpointRepository.GetSingleEndpoint(id)
}

func (e *EndpointService) GetAllEndpoints() ([]*models.Endpoint, error) {
	return e.EndpointRepository.GetAllEndpoints()
}

func (e *EndpointService) GetAllEndpointsWithLatestLog() ([]*models.Result, error) {
	return e.EndpointRepository.GetAllEndpointsWithLatestLog()
}

func (e *EndpointService) DeleteEndpoint(id string) {
	_, err := e.EndpointRepository.DeleteEndpoint(id)
	if err != nil {
		return
	}
}
