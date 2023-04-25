package handler

import (
	"github.com/aaltgod/bezdna/internal/service"
)

type handler struct {
	service service.Service
}

func New(
	service service.Service,
) Handler {
	return &handler{
		service: service,
	}
}
