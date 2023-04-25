package service

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/aaltgod/bezdna/internal/sniffer"
	log "github.com/sirupsen/logrus"
)

func (s *service) AddService(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Error(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	service := domain.Service{}

	if err = json.Unmarshal(body, &service); err != nil {
		log.Error(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	if err := s.dbRepository.InsertService(service); err != nil {
		log.Error(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	if err := s.sniffer.AddConfig(sniffer.Config{
		ServiceName: service.Name,
		Port:        service.Port,
	}); err != nil {
		log.Println(err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}
