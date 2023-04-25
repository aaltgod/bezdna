package db

import (
	"time"

	"github.com/aaltgod/bezdna/internal/domain"
)

type Service struct {
	Name string `db:"name"`
	Port uint16 `db:"port"`
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

type Stream struct {
	Ack       uint64    `db:"ack"`
	Timestamp time.Time `db:"timestamp"`
	Payload   string    `db:"payload"`
}

func (s Stream) ToDomain() domain.Stream {
	return domain.Stream{
		Ack:       s.Ack,
		Timestamp: s.Timestamp,
		Payload:   s.Payload,
	}
}

type Streams []Stream

func (s Streams) ToDomain() []domain.Stream {
	result := make([]domain.Stream, 0, len(s))

	for _, stream := range s {
		result = append(result, stream.ToDomain())
	}

	return result
}
