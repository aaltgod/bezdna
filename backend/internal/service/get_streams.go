package service

import (
	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/pkg/errors"
)

func (s *service) GetStreams(
	id, limit int64,
) ([]domain.Stream, error) {
	streams, err := s.dbRepository.GetStreams(id, limit)
	if err != nil {
		return nil, errors.Wrap(err, "dbRepository.GetStreams")
	}

	streamIDs := make([]int64, 0, len(streams))
	for _, stream := range streams {
		streamIDs = append(streamIDs, stream.ID)
	}

	flags, err := s.dbRepository.GetFlagsByStreamIDs(streamIDs)
	if err != nil {
		return nil, errors.Wrap(err, "dbRepository.GetFlagsByStreamIDs")
	}

	for i, stream := range streams {
		streams[i].Flags = flags[stream.ID]
	}

	return streams, nil
}

func (s *service) GetLastStreams(limit int64) ([]domain.Stream, error) {
	streams, err := s.dbRepository.GetLastStreams(limit)
	if err != nil {
		return nil, errors.Wrap(err, "dbRepository.GetLastStreams")
	}

	streamIDs := make([]int64, 0, len(streams))
	for _, stream := range streams {
		streamIDs = append(streamIDs, stream.ID)
	}

	flags, err := s.dbRepository.GetFlagsByStreamIDs(streamIDs)
	if err != nil {
		return nil, errors.Wrap(err, "dbRepository.GetFlagsByStreamIDs")
	}

	for i, stream := range streams {
		streams[i].Flags = flags[stream.ID]
	}

	return streams, nil
}
