package handler

import "net/http"

type Handler interface {
	CreateService(w http.ResponseWriter, req *http.Request)
	GetServices(w http.ResponseWriter, req *http.Request)

	GetStreamsByService(w http.ResponseWriter, req *http.Request)
	WSGetStreams(w http.ResponseWriter, req *http.Request)
}
