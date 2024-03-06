package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/zhayt/cert-tz/internal/model"
	service2 "github.com/zhayt/cert-tz/internal/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// HashCalculation calculates the hash for a given certificate.
//
// @Summary Calculate certificate hash
// @Description Calculates the hash for a given certificate.
// @Tags Hash
// @Accept json
// @Produce json
// @Param certHash body model.CertHash true "Certificate hash details"
// @Success 200 {object} struct{ID uint64 `json:"id"`} "Successful response"
// @Failure 400 {object} ErrorResponse "Not found"
// @Failure 503 {object} struct { Message string `json:"message"` } "Service unavailable"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /hash/calc [post]
func (h *Handler) HashCalculation(w http.ResponseWriter, r *http.Request) {
	var certHash model.CertHash

	err := json.NewDecoder(r.Body).Decode(&certHash)
	defer r.Body.Close()
	if err != nil {
		h.l.Error("Decode error", zap.Error(err))

		h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	certHash.ID, err = h.service.Hash.CalculateHash(certHash)
	if err != nil {
		h.l.Error("CalculateHash", zap.Error(err))

		if errors.Is(err, service2.ErrInvalidData) {
			h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		} else if errors.Is(err, service2.ErrWorkersPool) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(503)

			json.NewEncoder(w).Encode(struct {
				Message string `json:"message"`
			}{
				Message: "The maximum number of hashes that can be computed simultaneously has been reached, try again later",
			})
		} else {
			h.respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")

	h.l.Info("Start calculate hash: success", zap.Uint64("id", certHash.ID))
	json.NewEncoder(w).Encode(struct {
		ID uint64 `json:"id"`
	}{certHash.ID})
}

// ShowHash retrieves the calculated hash for a given ID.
//
// @Summary Get calculated hash
// @Description Retrieves the calculated hash for a given ID.
// @Tags Hash
// @Accept json
// @Produce json
// @Param id path integer true "Hash ID"
// @Success 200 {object} model.CertHash "Successful response"
// @Failure 404 {object} ErrorResponse "Not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /hash/result/{id} [get]
func (h *Handler) ShowHash(w http.ResponseWriter, r *http.Request) {
	hashID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || hashID <= 0 {
		h.l.Error("Param error", zap.Error(err))

		h.respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	certHash, err := h.service.Hash.GetCalculatedHash(uint64(hashID))
	if err != nil {
		h.l.Error("GetCalculatedHash error", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			h.respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		h.respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	h.l.Info("Calculated Hash found", zap.Int("id", hashID), zap.String("hash", certHash.Hash))
	json.NewEncoder(w).Encode(certHash)
}
