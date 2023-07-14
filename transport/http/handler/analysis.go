package handler

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type AnalysisRequest struct {
	Text string `json:"text"`
}

type AnalysisResponseEmail struct {
	Emails []string `json:"emails"`
}

func (h *Handler) AnalysisToEmail(w http.ResponseWriter, r *http.Request) {
	var request AnalysisRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()
	if err != nil {
		h.l.Error("Decode error", zap.Error(err))

		h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	response := AnalysisResponseEmail{}

	response.Emails = h.service.Analysis.FindEmails(request.Text)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
