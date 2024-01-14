package db

import (
	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

func (r *repository) UpsertService(service domain.Service) error {
	_, err := r.db.Pool.Exec(
		queryUpsertService,
		service.Name,
		service.Port,
		service.FlagRegexp,
	)
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

		if err = rows.Scan(
			&service.Name,
			&service.Port,
			&service.FlagRegexp,
		); err != nil {
			return nil, errors.Wrap(err, WrapScan)
		}

		result = append(result, service)
	}

	return result.ToDomain(), nil
}
