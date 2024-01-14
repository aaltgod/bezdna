package db

import (
	"time"

	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r *repository) InsertStreams(streams []domain.Stream) ([]int64, error) {
	streamAmount := len(streams)
	if streamAmount == 0 {
		return nil, nil
	}

	var (
		serviceNames = make([]string, 0, streamAmount)
		servicePorts = make([]int32, 0, streamAmount)
		texts        = make([]*string, 0, streamAmount)
		flagRegexps  = make([]string, 0, streamAmount)
		startedAt    = make([]time.Time, 0, streamAmount)
		endedAt      = make([]time.Time, 0, streamAmount)
	)

	for _, stream := range streams {
		serviceNames = append(serviceNames, stream.ServiceName)
		servicePorts = append(servicePorts, stream.ServicePort)
		texts = append(texts, stream.Text)
		flagRegexps = append(flagRegexps, stream.FlagRegexp)
		startedAt = append(startedAt, stream.StartedAt)
		endedAt = append(endedAt, stream.EndedAt)
	}

	rows, err := r.db.Pool.Query(
		queryInsertStreams,
		pq.Array(serviceNames),
		pq.Array(servicePorts),
		pq.Array(texts),
		pq.Array(flagRegexps),
		pq.Array(startedAt),
		pq.Array(endedAt),
	)
	if err != nil {
		return nil, errors.Wrap(err, WrapQuery)
	}

	defer rows.Close()

	result := make([]int64, 0, streamAmount)

	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			return nil, errors.Wrap(err, WrapScan)
		}

		result = append(result, id)
	}

	return result, nil
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

func (r *repository) GetLastStreams(limit int64) ([]domain.Stream, error) {
	rows, err := r.db.Pool.Query(
		queryGetLastStreams,
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
			&stream.FlagRegexp,
			&stream.StartedAt,
			&stream.EndedAt,
		); err != nil {
			return nil, errors.Wrap(err, WrapScan)
		}

		result = append(result, stream)
	}

	return result.ToDomain(), nil
}

func (r *repository) GetStreams(
	id, limit int64,
) ([]domain.Stream, error) {
	rows, err := r.db.Pool.Query(
		queryGetStreams,
		id,
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
			&stream.FlagRegexp,
			&stream.StartedAt,
			&stream.EndedAt,
		); err != nil {
			return nil, errors.Wrap(err, WrapScan)
		}

		result = append(result, stream)
	}

	return result.ToDomain(), nil
}
