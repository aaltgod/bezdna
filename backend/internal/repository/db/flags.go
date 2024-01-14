package db

import (
	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r *repository) InsertFlags(flags []domain.Flag) error {
	if len(flags) == 0 {
		return nil
	}

	var (
		streamIDs  = make([]int64, 0, len(flags))
		texts      = make([]string, 0, len(flags))
		directions = make([]string, 0, len(flags))
	)

	for _, flag := range flags {
		streamIDs = append(streamIDs, flag.StreamID)
		texts = append(texts, flag.Text)
		directions = append(directions, flag.Direction.String())
	}

	_, err := r.db.Pool.Exec(
		queryInsertFlags,
		pq.Array(streamIDs),
		pq.Array(texts),
		pq.Array(directions),
	)
	if err != nil {
		return errors.Wrap(err, WrapExec)
	}

	return nil
}

func (r *repository) GetFlagsByStreamIDs(ids []int64) (map[int64][]domain.Flag, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	rows, err := r.db.Pool.Query(
		queryGetFlagsByStreamIDs,
		pq.Array(ids),
	)
	if err != nil {
		return nil, errors.Wrap(err, WrapQuery)
	}

	defer rows.Close()

	result := make(map[int64][]domain.Flag)

	for rows.Next() {
		flag := flag{}

		if err := rows.Scan(
			&flag.ID,
			&flag.StreamID,
			&flag.Text,
			&flag.Direction,
		); err != nil {
			return nil, errors.Wrap(err, WrapScan)
		}

		result[flag.StreamID] = append(result[flag.StreamID], flag.ToDomain())
	}

	return result, nil
}
