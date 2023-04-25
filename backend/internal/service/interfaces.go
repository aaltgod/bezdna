package service

import (
	"github.com/aaltgod/bezdna/internal/domain"
)

type Service interface {
	AddService(service domain.Service) error
	GetServices() ([]domain.Service, error)

	GetStreamsByService(
		getStreamsByService domain.GetStreamsByService,
	) ([]domain.Stream, error)
}
