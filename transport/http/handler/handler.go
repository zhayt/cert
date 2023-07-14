package handler

import (
	"github.com/zhayt/cert-tz/service"
	"go.uber.org/zap"
	"time"
)

const _defaultContextTimeout = 5 * time.Second

type Handler struct {
	service *service.Service
	l       *zap.Logger
}

func NewHandler(service *service.Service, l *zap.Logger) *Handler {
	return &Handler{service: service, l: l}
}
