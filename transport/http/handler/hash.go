package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/zhayt/cert-tz/model"
	"github.com/zhayt/cert-tz/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

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

		if errors.Is(err, service.ErrInvalidData) {
			h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		} else if errors.Is(err, service.ErrWorkersPool) {
			w.Header().Set("Content-Type", "application/json")

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
