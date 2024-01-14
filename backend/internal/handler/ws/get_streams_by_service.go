package ws

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/aaltgod/bezdna/pkg/helpers"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	log "github.com/sirupsen/logrus"
)

type GetServicesByStreamRequest struct {
	ServiceName string `json:"service_name"`
	ServicePort int32  `json:"service_port"`
	Offset      int64  `json:"offset"`
	Limit       int64  `json:"limit"`
}

type GetServiceByStreamResponse struct {
	Streams []Stream `json:"streams"`
}

type Flag struct {
	Text      string `json:"text"`
	Direction string `json:"direction"`
}

type Stream struct {
	ID          int64  `json:"id"`
	ServiceName string `json:"service_name"`
	ServicePort int32  `json:"service_port"`
	Text        string `json:"text"`
	FlagRegexp  string `json:"flag_regexp"`
	Flags       []Flag `json:"flags"`
	StartedAt   string `json:"started_at"`
	EndedAt     string `json:"ended_at"`
}

func (h *handler) GetStreamsByService(w http.ResponseWriter, req *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(req, w)
	if err != nil {
		http.Error(w, "err", http.StatusInternalServerError)
		return
	}

	go func() {
		defer conn.Close()

		log.Info("CONNECTED")

		data, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Error("ReadClientText ", err)
			return
		}

		errCh := make(chan error)

		// handling client error
		go func() {
			_, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				errCh <- err
			}
		}()

		req := GetServicesByStreamRequest{}

		if err := json.Unmarshal(data, &req); err != nil {
			log.Error("Unmarshal ", err)
			return
		}

		var (
			ticker = time.NewTicker(time.Microsecond * 5)
			offset = req.Offset
		)

		for {
			select {
			case err = <-errCh:
				log.Error("Connection error ", err)
				return

			case <-ticker.C:
				streams, err := h.service.GetStreamsByService(
					domain.Service{
						Name: req.ServiceName,
						Port: req.ServicePort,
					},
					offset,
					req.Limit,
				)
				if err != nil {
					log.Error("service.GetStreamsByService ", err)
					return
				}

				respStreams := make([]Stream, 0, len(streams))

				for _, stream := range streams {
					respStreams = append(respStreams, Stream{
						ID:          stream.ID,
						ServiceName: stream.ServiceName,
						ServicePort: stream.ServicePort,
						FlagRegexp:  stream.FlagRegexp,
						Text:        *stream.Text,
						StartedAt:   stream.StartedAt.Format(time.TimeOnly),
						EndedAt:     stream.EndedAt.Format(time.TimeOnly),
					})
				}

				for _, batch := range helpers.Batch(respStreams, 5) {
					select {
					case err = <-errCh:
						log.Error("Connection error ", err)
						return

					default:
						res, err := json.Marshal(GetServiceByStreamResponse{
							Streams: batch,
						})
						if err != nil {
							log.Error("Marshal: ", err)
							return
						}

						err = wsutil.WriteServerMessage(conn, ws.OpCode(ws.StateServerSide), res)
						if err != nil {
							http.Error(w, "err", http.StatusInternalServerError)
							return
						}

						log.Info("SEND")

						time.Sleep(3 * time.Second)
					}
				}

				ticker.Reset(5 * time.Second)
			}
		}
	}()
}
