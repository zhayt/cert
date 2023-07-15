package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CounterResponse struct {
	Value string `json:"value"`
}

func (h *Handler) CounterIncrease(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeout)
	defer cancel()

	val, err := strconv.Atoi(mux.Vars(r)["i"])
	if err != nil {
		h.l.Error("update user: url param error", zap.Error(err))

		h.respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if err := h.service.Counter.IncreaseCounter(ctx, "counter", int64(val)); err != nil {
		h.l.Error("IncreaseCounter error", zap.Error(err))

		h.respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	h.l.Info("counter increased", zap.Int("val", val))
	h.respondWithCommonSuccess(w, "counter increased successfully")
}

func (h *Handler) CounterDecrease(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeout)
	defer cancel()

	val, err := strconv.Atoi(mux.Vars(r)["i"])
	if err != nil {
		h.l.Error("update user: url param error", zap.Error(err))

		h.respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	if err := h.service.Counter.DecreaseCounter(ctx, "counter", int64(val)); err != nil {
		h.l.Error("DecreaseCounter error", zap.Error(err))

		h.respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	h.l.Info("counter decreased", zap.Int("val", val))
	h.respondWithCommonSuccess(w, "counter decreased successfully")
}

func (h *Handler) ShowCounter(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeout)
	defer cancel()

	val, err := h.service.Counter.GetCounter(ctx, "counter")
	if err != nil {
		h.l.Error("GetCounter error", zap.Error(err))

		val = "0"
	}

	response := CounterResponse{Value: val}

	h.l.Info("counter value got", zap.String("val", val))
	json.NewEncoder(w).Encode(response)
}
