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
	GetServiceByPort(port int32) (*domain.Service, error)
	GetServices() ([]domain.Service, error)
}
type Streamer interface {
	InsertStreams(streams []domain.Stream) error
	GetStreamsByService(
		service domain.Service, offset, limit int64,
	) ([]domain.Stream, error)
}
