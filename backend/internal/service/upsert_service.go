package service

import (
	"regexp"

	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/aaltgod/bezdna/internal/sniffer"
	"github.com/pkg/errors"
)

func (s *service) UpsertService(service domain.Service) error {
	flagRegexp, err := regexp.Compile(service.FlagRegexp)
	if err != nil {
		return errors.Wrap(err, "Compile")
	}

	if err := s.dbRepository.UpsertService(service); err != nil {
		return errors.Wrap(err, WrapUpsertService)
	}

	if err := s.sniffer.AddConfig(sniffer.Config{
		ServiceName: service.Name,
		ServicePort: service.Port,
		FlagRegexp:  flagRegexp,
	}); err != nil {
		return errors.Wrap(err, WrapAddConfig)
	}

	return nil
}
