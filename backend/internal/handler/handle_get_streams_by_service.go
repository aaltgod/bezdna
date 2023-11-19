package handler

import (
	"encoding/json"
	"github.com/aaltgod/bezdna/internal/domain"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (h *handler) GetStreamsByService(w http.ResponseWriter, req *http.Request) {
	_, err := io.ReadAll(req.Body)
	if err != nil {
		log.Errorln(WrapfGetStreamsByService(err, WrapReadAll))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return
	}

	streams := []domain.Fake{
		{
			UserId:    1,
			Id:        12121212,
			Completed: true,
			Title:     "TEXT",
		},
		{
			UserId:    2,
			Id:        12121212,
			Completed: true,
			Title:     "qwer",
		},
		{
			UserId:    1,
			Id:        12121212,
			Completed: true,
			Title:     "yrt",
		},
		{
			UserId:    2,
			Id:        12121212,
			Completed: true,
			Title:     "ert",
		},
		{
			UserId:    1,
			Id:        12121212,
			Completed: true,
			Title:     "wer",
		},
		{
			UserId:    2,
			Id:        12121212,
			Completed: true,
			Title:     "r",
		},
		{
			UserId:    1,
			Id:        12121212,
			Completed: true,
			Title:     "wer",
		},
		{
			UserId:    2,
			Id:        12121212,
			Completed: true,
			Title:     "wer",
		},
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

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)

	return

	// getStreamsByService := domain.GetStreamsByService{}

	// if err = json.Unmarshal(body, &getStreamsByService); err != nil {
	// 	log.Errorln(WrapfGetStreamsByService(err, WrapUnmarshal))

	// 	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

	// 	return
	// }

	// // FIXME: offset can be zero
	// if getStreamsByService.Offset < 1 {
	// 	log.Errorln(ErrMinOffset, getStreamsByService)

	// 	http.Error(w, ErrMinOffset.Error(), http.StatusBadRequest)

	// 	return
	// }

	// if getStreamsByService.Offset > 20 {
	// 	log.Errorln(ErrMaxOffset, getStreamsByService)

	// 	http.Error(w, ErrMaxOffset.Error(), http.StatusBadRequest)

	// 	return
	// }

	// if getStreamsByService.Limit > 20 {
	// 	log.Errorln(ErrMaxLimit, getStreamsByService)

	// 	http.Error(w, ErrMaxLimit.Error(), http.StatusBadRequest)

	// 	return
	// }

	// if getStreamsByService.Limit == 0 {
	// 	getStreamsByService.Limit = 20
	// }

	// streams, err := h.service.GetStreamsByService(getStreamsByService)
	// if err != nil {
	// 	log.Errorln(WrapfGetStreamsByService(err, WrapGetStreamsByService))

	// 	http.Error(
	// 		w,
	// 		http.StatusText(http.StatusInternalServerError),
	// 		http.StatusInternalServerError,
	// 	)

	// 	return
	// }

	// result, err := json.Marshal(streams)
	// if err != nil {
	// 	log.Errorln(WrapfGetStreamsByService(err, WrapMarshal))

	// 	http.Error(
	// 		w,
	// 		http.StatusText(http.StatusInternalServerError),
	// 		http.StatusInternalServerError,
	// 	)

	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// w.Write(result)
}
