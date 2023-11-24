package service

import (
	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/pkg/errors"
)

func (s *service) GetStreamsByService(
	service domain.Service, offset, limit int64,
) ([]domain.Stream, error) {
	streams, err := s.dbRepository.GetStreamsByService(
		service,
		offset,
		limit,
	)
	if err != nil {
		return nil, errors.Wrap(err, WrapGetStreamsByService)
	}

	return streams, nil
}
