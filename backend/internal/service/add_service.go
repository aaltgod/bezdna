package service

import (
	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/aaltgod/bezdna/internal/sniffer"
	"github.com/pkg/errors"
)

func (s *service) AddService(service domain.Service) error {
	if err := s.dbRepository.InsertService(service); err != nil {
		return errors.Wrap(err, WrapInsertService)
	}

	if err := s.sniffer.AddConfig(sniffer.Config{
		ServiceName: service.Name,
		Port:        service.Port,
	}); err != nil {
		return errors.Wrap(err, WrapAddConfig)
	}

	return nil
}
