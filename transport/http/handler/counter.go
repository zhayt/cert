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

// CounterIncrease increases the counter value by the specified amount.
//
// @Summary Increase the counter value
// @Description Increases the counter value by the specified amount.
// @Tags Counter
// @Accept json
// @Produce json
// @Param i path integer true "Amount to increase the counter value by"
// @Success 200 {object} CounterResponse "Successful response"
// @Failure 400 {object} ErrorResponse "Not Found"
// @Failure 500 {object} ErrorResponse "Internal Error"
// @Router /counter/add/{i} [post]
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

// CounterDecrease decreases the counter value by the specified amount.
//
// @Summary Decrease the counter value
// @Description Decreases the counter value by the specified amount.
// @Tags Counter
// @Accept json
// @Produce json
// @Param i path integer true "Amount to decrease the counter value by"
// @Success 200 {object} CounterResponse "Successful response"
// @Failure 400 {object} ErrorResponse "Not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /counter/decrease/{i} [post]
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

// ShowCounter retrieves the current value of the counter.
//
// @Summary Get the counter value
// @Description Retrieves the current value of the counter.
// @Tags Counter
// @Accept json
// @Produce json
// @Success 200 {object} CounterResponse "Successful response"
// @Router /counter/val [get]
func (h *Handler) ShowCounter(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeout)
	defer cancel()

	val, err := h.service.Counter.GetCounter(ctx, "counter")
	if err != nil {
		h.l.Error("GetCounter error", zap.Error(err))

		val = "0"
	}

	response := CounterResponse{Value: val}

	w.Header().Set("Content-Type", "application/json")

	h.l.Info("counter value got", zap.String("val", val))
	json.NewEncoder(w).Encode(response)
}
