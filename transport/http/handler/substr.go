package handler

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type Input struct {
	Str string `json:"input"`
}

type Output struct {
	Str string `json:"output"`
}

func (h *Handler) FindSubString(w http.ResponseWriter, r *http.Request) {
	var input Input

	err := json.NewDecoder(r.Body).Decode(&input)
	defer r.Body.Close()
	if err != nil {
		h.l.Error("Decode error", zap.Error(err))

		h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	output := Output{}
	output.Str, err = h.service.SubStr.LongestSubstring(input.Str)
	if err != nil {
		h.l.Error("LongestSubstring error", zap.Error(err))

		h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	h.l.Info("sub string found", zap.String("substring", output.Str))
	json.NewEncoder(w).Encode(output)
}
