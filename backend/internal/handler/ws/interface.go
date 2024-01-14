package ws

import "net/http"

type Handler interface {
	GetStreams(w http.ResponseWriter, req *http.Request)
	GetStreamsByService(w http.ResponseWriter, req *http.Request)
}
