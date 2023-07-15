package handler

import (
	"bytes"
	"encoding/json"
	"github.com/zhayt/cert-tz/service"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_FindSubString1(t *testing.T) {
	type fields struct {
		service *service.Service
		l       *zap.Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	field := fields{
		service: service.NewService(nil, zap.NewExample()),
		l:       zap.NewNop(), // Use a no-op logger for testing
	}

	resRecorder := httptest.NewRecorder()
	reqBody := bytes.NewBufferString(`{"input": "abcdabcde"}`)
	req, err := http.NewRequest("POST", "/find-substring", reqBody)

	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"success", field, args{resRecorder, req}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				service: tt.fields.service,
				l:       tt.fields.l,
			}
			h.FindSubString(tt.args.w, tt.args.r)
			var output Output
			if err := json.Unmarshal(resRecorder.Body.Bytes(), &output); err != nil {
				t.Fatalf("Failed to parse response body: %v", err)
			}

			// Check the expected output
			expectedOutput := Output{
				Str: "abcde",
			}
			if !reflect.DeepEqual(output, expectedOutput) {
				t.Errorf("Expected output %v, got %v", expectedOutput, output)
			}
		})
	}
}
