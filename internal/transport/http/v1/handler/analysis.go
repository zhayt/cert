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

type AnalysisResponseIINs struct {
	IINs []string `json:"iins"`
}

// AnalysisToEmail analyzes the given text to extract emails.
//
// @Summary Analyze text to extract emails
// @Description Analyzes the provided text to identify and extract emails.
// @Tags Analysis
// @Accept json
// @Produce json
// @Param request body AnalysisRequest true "Text to analyze"
// @Success 200 {object} AnalysisResponseEmail "Successful response"
// @Failure 400 {object} ErrorResponse "Not found"
// @Router /email/check [post]
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

	h.l.Info("Text analyzed to emails", zap.Int("count", len(response.Emails)))
	json.NewEncoder(w).Encode(response)
}

// AnalysisToIIN analyzes the given text to extract IINs.
//
// @Summary Analyze text to extract IINs
// @Description Analyzes the provided text to identify and extract IINs (Individual Identification Numbers).
// @Tags Analysis
// @Accept json
// @Produce json
// @Param request body AnalysisRequest true "Text to analyze"
// @Success 200 {object} AnalysisResponseIINs "Successful response"
// @Failure 400 {object} ErrorResponse "Not found"
// @Router /iin/check [post]
func (h *Handler) AnalysisToIIN(w http.ResponseWriter, r *http.Request) {
	var request AnalysisRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()
	if err != nil {
		h.l.Error("Decode error", zap.Error(err))

		h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	response := AnalysisResponseIINs{}

	h.l.Info(request.Text)
	response.IINs = h.service.Analysis.FindsIINs(request.Text)

	w.Header().Set("Content-Type", "application/json")

	h.l.Info("Text analyzed to iins", zap.Int("count", len(response.IINs)))
	json.NewEncoder(w).Encode(response)
}
