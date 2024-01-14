package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (h *handler) GetServices(w http.ResponseWriter, req *http.Request) {
	services, err := h.service.GetServices()
	if err != nil {
		log.Errorln(WrapfGetServices(err, WrapCreateService))

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}

	response, err := json.Marshal(services)
	if err != nil {
		log.Errorln(WrapfGetServices(err, WrapMarshal))

		http.Error(
			w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
