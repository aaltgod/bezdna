package service

import "net/http"

type Service interface {
	AddService(w http.ResponseWriter, req *http.Request)
	GetServices(w http.ResponseWriter, req *http.Request)
}
