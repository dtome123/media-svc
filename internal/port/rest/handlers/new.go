package handlers

import (
	"media-svc/internal/services"
)

type impl struct {
	svc *services.Service
}

func NewHandler(svc *services.Service) Handler {
	return &impl{
		svc: svc,
	}
}
