package db

import (
	"github.com/aaltgod/bezdna/internal/domain"
)

type Repository interface {
	Servicer
	Streamer
	Flager
}

type Servicer interface {
	UpsertService(service domain.Service) error
	GetServiceByPort(port int32) (*domain.Service, error)
	GetServices() ([]domain.Service, error)
}

type Streamer interface {
	InsertStreams(streams []domain.Stream) ([]int64, error)
	GetStreams(
		id, limit int64,
	) ([]domain.Stream, error)
	GetLastStreams(limit int64) ([]domain.Stream, error)
	GetStreamsByService(
		service domain.Service, offset, limit int64,
	) ([]domain.Stream, error)
}

type Flager interface {
	InsertFlags(flags []domain.Flag) error
	GetFlagsByStreamIDs(ids []int64) (map[int64][]domain.Flag, error)
}
