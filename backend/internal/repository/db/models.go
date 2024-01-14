package db

import (
	"time"

	"github.com/aaltgod/bezdna/internal/domain"
)

type Service struct {
	Name       string `db:"name"`
	Port       int32  `db:"port"`
	FlagRegexp string `db:"flag_regexp"`
}

func (s Service) ToDomain() domain.Service {
	return domain.Service{
		Name:       s.Name,
		Port:       s.Port,
		FlagRegexp: s.FlagRegexp,
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
	FlagRegexp  string    `db:"flag_regexp"`
	FlagAction  string    `db:"flag_action"`
	StartedAt   time.Time `db:"started_at"`
	EndedAt     time.Time `db:"ended_at"`
}

func (s stream) ToDomain() domain.Stream {
	return domain.Stream{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		ServicePort: s.ServicePort,
		Text:        s.Text,
		FlagRegexp:  s.FlagRegexp,
		FlagAction:  domain.FlagDirection(s.FlagAction),
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

type flag struct {
	ID        int64  `db:"id"`
	StreamID  int64  `db:"stream_id"`
	Text      string `db:"text"`
	Direction string `db:"direction"`
}

func (f flag) ToDomain() domain.Flag {
	return domain.Flag{
		ID:        f.ID,
		StreamID:  f.StreamID,
		Text:      f.Text,
		Direction: domain.FlagDirection(f.Direction),
	}
}

type flags []flag

func (f flags) ToDomain() []domain.Flag {
	result := make([]domain.Flag, 0, len(f))

	for _, flag := range f {
		result = append(result, flag.ToDomain())
	}

	return result
}
