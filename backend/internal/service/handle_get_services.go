package service

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *service) GetServices(w http.ResponseWriter, req *http.Request) {
	services, err := s.dbRepository.GetServices()
	if err != nil {
		log.Error(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(services)
	if err != nil {
		log.Error(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
