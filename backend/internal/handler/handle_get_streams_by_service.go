package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aaltgod/bezdna/internal/domain"
	log "github.com/sirupsen/logrus"
)

func (h *handler) GetStreamsByService(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Errorln(WrapfGetStreamsByService(err, WrapReadAll))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	getStreamsByService := domain.GetStreamsByService{}

	if err = json.Unmarshal(body, &getStreamsByService); err != nil {
		log.Errorln(WrapfGetStreamsByService(err, WrapUnmarshal))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	streams, err := h.service.GetStreamsByService(getStreamsByService)
	if err != nil {
		log.Errorln(WrapfGetStreamsByService(err, WrapGetStreamsByService))

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}

	result, err := json.Marshal(streams)
	if err != nil {
		log.Errorln(WrapfGetStreamsByService(err, WrapMarshal))

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
