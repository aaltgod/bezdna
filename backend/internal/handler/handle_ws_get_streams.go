package handler

import (
	"net/http"
	"time"

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

		// err = wsutil.WriteServerMessage(conn, ws.OpCode(ws.StateServerSide), []byte("STREAM"))
		// if err != nil {
		// 	http.Error(w, "err", http.StatusInternalServerError)

		// 	return
		// }

		for {
			err = wsutil.WriteServerMessage(conn, ws.OpCode(ws.StateServerSide), []byte("STREAM"))
			if err != nil {
				http.Error(w, "err", http.StatusInternalServerError)

				return
			}

			log.Info("SEND")

			time.Sleep(3 * time.Second)
		}
	}()
}
