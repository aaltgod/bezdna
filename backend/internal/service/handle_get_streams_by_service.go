package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aaltgod/bezdna/internal/domain"
	log "github.com/sirupsen/logrus"
)

func (s *service) GetStreamsByService(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Error(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	getStreamsByService := domain.GetStreamsByService{}

	if err = json.Unmarshal(body, &getStreamsByService); err != nil {
		log.Error(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	streams, err := s.dbRepository.GetStreamsByService(domain.Service{
		Name: getStreamsByService.Name,
		Port: getStreamsByService.Port,
	},
		getStreamsByService.Offset,
	)
	if err != nil {
		log.Error(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	result, err := json.Marshal(streams)
	if err != nil {
		log.Error(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
