package db

import (
	"time"

	"github.com/aaltgod/bezdna/internal/database"
	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/jackc/pgx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type repository struct {
	db *database.DBAdapter
}

func New(db *database.DBAdapter) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) InsertService(service domain.Service) error {
	_, err := r.db.Pool.Exec(queryInsertService, service.Name, service.Port)
	if err != nil {
		return errors.Wrap(err, WrapExec)
	}

	return nil
}

func (r *repository) GetServiceByPort(port int32) (*domain.Service, error) {
	service := Service{}

	if err := r.db.Pool.QueryRow(
		queryGetServiceByPort,
		port,
	).Scan(
		&service.Name,
		&service.Port,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, WrapScan)
	}

	domainService := service.ToDomain()

	return &domainService, nil
}

func (r *repository) GetServices() ([]domain.Service, error) {
	rows, err := r.db.Pool.Query(queryGetServices)
	if err != nil {
		return nil, errors.Wrap(err, WrapQuery)
	}
	defer rows.Close()

	var result Services

	for rows.Next() {
		service := Service{}

		if err = rows.Scan(&service.Name, &service.Port); err != nil {
			return nil, errors.Wrap(err, WrapScan)
		}

		result = append(result, service)
	}

	return result.ToDomain(), nil
}

func (r *repository) InsertStreams(streams []domain.Stream) error {
	streamAmount := len(streams)
	if streamAmount == 0 {
		return nil
	}

	var (
		serviceNames = make([]string, 0, streamAmount)
		servicePorts = make([]int32, 0, streamAmount)
		texts        = make([]*string, 0, streamAmount)
		startedAt    = make([]time.Time, 0, streamAmount)
		endedAt      = make([]time.Time, 0, streamAmount)
	)

	for _, stream := range streams {
		serviceNames = append(serviceNames, stream.ServiceName)
		servicePorts = append(servicePorts, stream.ServicePort)
		texts = append(texts, stream.Text)
		startedAt = append(startedAt, stream.StartedAt)
		endedAt = append(endedAt, stream.EndedAt)
	}

	if _, err := r.db.Pool.Exec(
		queryInsertStreams,
		pq.Array(serviceNames),
		pq.Array(servicePorts),
		pq.Array(texts),
		pq.Array(startedAt),
		pq.Array(endedAt),
	); err != nil {
		return errors.Wrap(err, "db.Pool.Exec")
	}

	return nil
}

func (r *repository) GetStreamsByService(
	service domain.Service, offset, limit int64,
) ([]domain.Stream, error) {
	rows, err := r.db.Pool.Query(
		queryGetStreamsByService,
		service.Name,
		service.Port,
		offset,
		limit,
	)
	if err != nil {
		return nil, errors.Wrap(err, WrapQuery)
	}
	defer rows.Close()

	var result Streams

	for rows.Next() {
		stream := stream{}

		if err = rows.Scan(
			&stream.ID,
			&stream.ServiceName,
			&stream.ServicePort,
			&stream.Text,
			&stream.StartedAt,
			&stream.EndedAt,
		); err != nil {
			return nil, errors.Wrap(err, WrapScan)
		}

		result = append(result, stream)
	}

	return result.ToDomain(), nil
}
