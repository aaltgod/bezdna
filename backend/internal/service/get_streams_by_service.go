package service

import (
	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/pkg/errors"
)

func (s *service) GetStreamsByService(
	getStreamsByService domain.GetStreamsByService,
) ([]domain.Stream, error) {
	streams, err := s.dbRepository.GetStreamsByService(getStreamsByService)
	if err != nil {
		return nil, errors.Wrap(err, WrapGetStreamsByService)
	}

	return streams, nil
}
