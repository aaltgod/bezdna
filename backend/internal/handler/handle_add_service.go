package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aaltgod/bezdna/internal/domain"
	log "github.com/sirupsen/logrus"
)

func (h *handler) AddService(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Errorln(WrapfAddService(err, WrapReadAll))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	service := domain.Service{}

	if err = json.Unmarshal(body, &service); err != nil {
		log.Errorln(WrapfAddService(err, WrapUnmarshal))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	if err := h.service.AddService(service); err != nil {
		log.Errorln(WrapfAddService(err, WrapAddService))

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
