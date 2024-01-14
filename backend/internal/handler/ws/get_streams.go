package ws

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/aaltgod/bezdna/internal/domain"
	"github.com/aaltgod/bezdna/pkg/helpers"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	log "github.com/sirupsen/logrus"
)

type GetStreamsResponse struct {
	Offset  int64    `json:"offset"`
	Streams []Stream `json:"streams"`
}

func (h *handler) GetStreams(w http.ResponseWriter, req *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(req, w)
	if err != nil {
		http.Error(w, "err", http.StatusInternalServerError)
		return
	}

	go func() {
		defer conn.Close()

		log.Warn("CONNECTED GetStreams")

		errCh := make(chan error)

		// handling client error
		go func() {
			_, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				errCh <- err
			}
		}()

		var (
			limit      int64 = 100
			offset, id int64
			ticker     = time.NewTicker(time.Microsecond * 5)

			initilized bool
		)

		for {
			select {
			case err = <-errCh:
				log.Error("Connection error ", err)
				return

			case <-ticker.C:
				var streams []domain.Stream

				if !initilized {
					streams, err = h.service.GetLastStreams(limit)
					if err != nil {
						http.Error(w, "err", http.StatusInternalServerError)

						log.Error("service.GetLastStreams ", err)

						return
					}

					sort.Slice(streams, func(i, j int) bool {
						return streams[i].ID < streams[j].ID
					})

					initilized = true
				} else {
					streams, err = h.service.GetStreams(id, limit)
					if err != nil {
						log.Error("service.GetStreams ", err)
						return
					}
				}

				if len(streams) != 0 {
					id = streams[len(streams)-1].ID
				}

				respStreams := make([]Stream, 0, len(streams))

				for _, stream := range streams {
					flags := make([]Flag, 0, len(stream.Flags))

					for _, flag := range stream.Flags {
						flags = append(flags, Flag{
							Text:      flag.Text,
							Direction: flag.Direction.String(),
						})
					}
					respStreams = append(respStreams, Stream{
						ID:          stream.ID,
						ServiceName: stream.ServiceName,
						ServicePort: stream.ServicePort,
						FlagRegexp:  stream.FlagRegexp,
						Flags:       flags,
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
						// offset += int64(len(batch))

						res, err := json.Marshal(GetStreamsResponse{
							Streams: batch,
							Offset:  offset,
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

						log.Warn("SEND")

						time.Sleep(1 * time.Second)
					}
				}

				ticker.Reset(5 * time.Second)
			}
		}
	}()
}
