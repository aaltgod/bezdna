package service

import (
	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/pkg/errors"
)

func (s *service) GetServices() ([]domain.Service, error) {
	services, err := s.dbRepository.GetServices()
	if err != nil {
		return nil, errors.Wrap(err, WrapGetServices)
	}

	return services, nil
}
