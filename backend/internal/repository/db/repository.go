package db

import "github.com/aaltgod/bezdna/internal/database"

type repository struct {
	db *database.DBAdapter
}

func New(db *database.DBAdapter) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) InsertService(serviceName string, port uint16) error {
	_, err := r.db.Pool.Exec(queryInsertService, serviceName, port)
	if err != nil {
		return err
	}

	return nil
}
