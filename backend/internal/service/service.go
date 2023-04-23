package service

import (
	"github.com/aaltgod/bezdna/internal/repository/db"
	"github.com/aaltgod/bezdna/internal/sniffer"
)

type service struct {
	sniffer      *sniffer.Sniffer
	dbRepository db.Repository
}

func New(sniffer *sniffer.Sniffer, dbRepository db.Repository) Service {
	return &service{
		sniffer:      sniffer,
		dbRepository: dbRepository,
	}
}
