package db

import (
	"github.com/aaltgod/bezdna/internal/domain"
)

type Repository interface {
	Servicer
	Streamer
}

type Servicer interface {
	InsertService(service domain.Service) error
	GetServices() ([]domain.Service, error)
}
type Streamer interface {
	InsertStreamByService(
		stream domain.Stream, service domain.Service,
	) error
	GetStreamsByService(
		getStreamsByService domain.GetStreamsByService,
	) ([]domain.Stream, error)
}
