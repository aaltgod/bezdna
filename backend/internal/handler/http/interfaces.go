package handler

import "net/http"

type Handler interface {
	UpsertService(w http.ResponseWriter, req *http.Request)
	GetServices(w http.ResponseWriter, req *http.Request)

	GetStreamsByService(w http.ResponseWriter, req *http.Request)
}
