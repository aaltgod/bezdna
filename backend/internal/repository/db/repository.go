package db

import (
	"github.com/aaltgod/bezdna/internal/database"
	"github.com/aaltgod/bezdna/internal/domain"
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

func (r *repository) UpsertFlagRegexp(flagRegexp, serviceName string, servicePort int32) error {
	if _, err := r.db.Pool.Exec(
		queryUpsertFlagRegexp,
		flagRegexp,
		serviceName,
		servicePort,
	); err != nil {
		return errors.Wrap(err, "db.Pool.Exec")
	}

	return nil
}

func (r *repository) GetFlagRegexpsByServices(services []domain.Service) ([]string, error) {
	serviceAmount := len(services)
	if serviceAmount == 0 {
		return nil, nil
	}

	var (
		serviceNames = make([]string, 0, serviceAmount)
		servicePorts = make([]int32, 0, serviceAmount)
	)

	for _, service := range services {
		serviceNames = append(serviceNames, service.Name)
		servicePorts = append(servicePorts, service.Port)
	}

	rows, err := r.db.Pool.Query(
		queryGetRegexpsByServices,
		pq.Array(serviceNames),
		pq.Array(servicePorts),
	)
	if err != nil {
		return nil, errors.Wrap(err, WrapQuery)
	}
	defer rows.Close()

	result := make([]string, 0, len(services))

	for rows.Next() {
		var flagRegexp string

		if err = rows.Scan(
			&flagRegexp,
		); err != nil {
			return nil, errors.Wrap(err, WrapScan)
		}

		result = append(result, flagRegexp)
	}

	return result, nil
}
