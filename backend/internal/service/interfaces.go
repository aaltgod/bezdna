package service

import (
	"github.com/aaltgod/bezdna/internal/domain"
)

type Service interface {
	Servicer
	Streamer
}

type Servicer interface {
	UpsertService(service domain.Service) error
	GetServices() ([]domain.Service, error)
}

type Streamer interface {
	GetStreamsByService(
		service domain.Service, offset, limit int64,
	) ([]domain.Stream, error)
	GetStreams(
		id, limit int64,
	) ([]domain.Stream, error)
	GetLastStreams(limit int64) ([]domain.Stream, error)
}
