package db

import (
	"github.com/aaltgod/bezdna/internal/database"
	"github.com/aaltgod/bezdna/internal/domain"
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

func (r *repository) InsertStreamByService(
	stream domain.Stream, service domain.Service,
) error {
	_, err := r.db.Pool.Exec(
		queryInsertStream, service.Name, service.Port,
		stream.Ack, stream.Timestamp, stream.Payload)
	if err != nil {
		return errors.Wrap(err, WrapExec)
	}

	return nil
}

func (r *repository) GetStreamsByService(
	service domain.Service, offset int64,
) ([]domain.Stream, error) {
	rows, err := r.db.Pool.Query(
		queryGetStreamsByService,
		service.Name, service.Port, offset)
	if err != nil {
		return nil, errors.Wrap(err, WrapQuery)
	}
	defer rows.Close()

	var result Streams

	for rows.Next() {
		stream := Stream{}

		if err = rows.Scan(&stream.Ack, &stream.Timestamp, &stream.Payload); err != nil {
			return nil, errors.Wrap(err, WrapScan)
		}

		result = append(result, stream)
	}

	return result.ToDomain(), nil
}
