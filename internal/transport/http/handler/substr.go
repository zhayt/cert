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

// FindSubString finds the longest without repeating characters substring in the given input string.
//
// @Summary Find longest without repeating characters common substring
// @Description Finds the longest common substring in the given input string.
// @Tags Substring
// @Accept json
// @Produce json
// @Param input body Input true "Input string"
// @Success 200 {object} Output "Successful response"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Router /substr/find [post]
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
