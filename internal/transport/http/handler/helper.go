package handler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessUserResponse struct {
	UserID  uint64 `json:"user_id"`
	Message string `json:"message"`
}

type CommonSuccessResponse struct {
	Message string `json:"message"`
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) respondWithSuccess(w http.ResponseWriter, userID uint64, message string) {
	response := SuccessUserResponse{
		UserID:  userID,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) respondWithCommonSuccess(w http.ResponseWriter, message string) {
	response := CommonSuccessResponse{
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
