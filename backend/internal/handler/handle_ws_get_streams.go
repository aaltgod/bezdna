package handler

import (
	"net/http"
	"time"

	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	log "github.com/sirupsen/logrus"
)

func (h *handler) WSGetStreams(w http.ResponseWriter, req *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(req, w)
	if err != nil {
		http.Error(w, "err", http.StatusInternalServerError)

		return
	}

	go func() {
		defer conn.Close()

		log.Info("CONNECTED")

		for {
			data, err := wsutil.ReadClientText(conn)
			if err != nil {
				log.Error("ReadFrame ", err)
				return
			}

			log.Info(string(data))

			streams, err := h.service.GetStreamsByService(
				domain.Service{
					Name: "щипитули",
					Port: 8973,
				},
				0,
				1000,
			)
			if err != nil {
				log.Error("service.GetStreamsByService ", err)
				return
			}

			for _, stream := range streams {
				err = wsutil.WriteServerMessage(conn, ws.OpCode(ws.StateServerSide), []byte(*stream.Text))
				if err != nil {
					http.Error(w, "err", http.StatusInternalServerError)

					return
				}

				log.Info("SEND")

				time.Sleep(3 * time.Second)
			}
		}
	}()
}
