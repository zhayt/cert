package handler

import (
	"github.com/zhayt/cert-tz/internal/service"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
)

func TestHandler_AnalysisToEmail(t *testing.T) {
	// Create a new request with a sample JSON payload
	reqBody := bytes.NewBufferString(`{"text": "Please contact us at Email: new@example.org Email: 74@gmailcom"}`)
	req, err := http.NewRequest("POST", "/analysis-to-email", reqBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to capture the response
	resRecorder := httptest.NewRecorder()

	// Create a new instance of the Handler
	handler := &Handler{
		service: service.NewService(nil, zap.NewExample()),
		l:       zap.NewNop(), // Use a no-op logger for testing
	}

	// Call the AnalysisToEmail method
	handler.AnalysisToEmail(resRecorder, req)

	// Check the response status code
	if resRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resRecorder.Code)
	}

	// Parse the response body
	var response AnalysisResponseEmail
	if err := json.Unmarshal(resRecorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Check the expected output
	expectedResponse := AnalysisResponseEmail{
		Emails: []string{"new@example.org"},
	}
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Errorf("Expected response %v, got %v", expectedResponse, response)
	}
}
