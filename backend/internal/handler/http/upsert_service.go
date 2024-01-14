package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/aaltgod/bezdna/internal/domain"
	serv "github.com/aaltgod/bezdna/internal/service"
	log "github.com/sirupsen/logrus"
)

func (h *handler) UpsertService(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Errorln(WrapfCreateService(err, WrapReadAll))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	service := domain.Service{}

	if err = json.Unmarshal(body, &service); err != nil {
		log.Errorln(WrapfCreateService(err, WrapUnmarshal))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	if err := h.service.UpsertService(service); err != nil {
		log.Errorln(WrapfCreateService(err, WrapCreateService))

		if errors.Is(err, serv.ErrAlreadyExist) {
			http.Error(
				w,
				http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest,
			)

			return
		}

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
