package service

import (
	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/aaltgod/bezdna/internal/sniffer"
	"github.com/pkg/errors"
)

func (s *service) CreateService(service domain.Service) error {
	serviceFromDB, err := s.dbRepository.GetServiceByPort(service.Port)
	if err != nil {
		return errors.Wrap(err, "dbRepository.GetServices")
	}

	if serviceFromDB != nil {
		return ErrAlreadyExist
	}

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
