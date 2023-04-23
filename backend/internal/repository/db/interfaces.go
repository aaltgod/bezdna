package db

type Repository interface {
	InsertService(serviceName string, port uint16) error
}
