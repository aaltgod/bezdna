package service

import (
	"net/http"
)

func (s *service) GetServices(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("services"))
}
