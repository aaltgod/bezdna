package db

import (
	"time"

	"github.com/aaltgod/bezdna/internal/domain"
)

type Service struct {
	Name string `db:"name"`
	Port int32  `db:"port"`
}

func (s Service) ToDomain() domain.Service {
	return domain.Service{
		Name: s.Name,
		Port: s.Port,
	}
}

type Services []Service

func (s Services) ToDomain() []domain.Service {
	result := make([]domain.Service, 0, len(s))

	for _, service := range s {
		result = append(result, service.ToDomain())
	}

	return result
}

type stream struct {
	ID          int64     `db:"id"`
	ServiceName string    `db:"service_name"`
	ServicePort int32     `db:"service_port"`
	Text        *string   `db:"text"`
	StartedAt   time.Time `db:"started_at"`
	EndedAt     time.Time `db:"ended_at"`
}

func (s stream) ToDomain() domain.Stream {
	return domain.Stream{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		ServicePort: s.ServicePort,
		Text:        s.Text,
		StartedAt:   s.StartedAt,
		EndedAt:     s.EndedAt,
	}
}

type Streams []stream

func (s Streams) ToDomain() []domain.Stream {
	result := make([]domain.Stream, 0, len(s))

	for _, stream := range s {
		result = append(result, stream.ToDomain())
	}

	return result
}
