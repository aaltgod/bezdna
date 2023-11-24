package service

import (
	"github.com/aaltgod/bezdna/internal/domain"
)

type Service interface {
	CreateService(service domain.Service) error
	GetServices() ([]domain.Service, error)

	GetStreamsByService(
		service domain.Service, offset, limit int64,
	) ([]domain.Stream, error)
}
