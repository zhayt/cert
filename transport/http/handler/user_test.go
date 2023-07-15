package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zhayt/cert-tz/service"
	"github.com/zhayt/cert-tz/storage"
	"github.com/zhayt/cert-tz/storage/postgre"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_CreateUser(t *testing.T) {
	// Create PostgreSQL container request
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "cert_db",
			"POSTGRES_PASSWORD": "qwerty",
			"POSTGRES_USER":     "cert",
		},
	}

	// Start PostgreSQL container
	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		t.Error(err)
	}

	// Get host and port of PostgreSQL container
	port, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		t.Error(err)
	}

	host, err := dbContainer.Host(context.Background())
	if err != nil {
		t.Error(err)
	}

	// Create db connection string and connect
	dbURI := fmt.Sprintf("postgres://cert:qwerty@%v:%v/cert_db", host, port.Port())

	db, err := sqlx.Connect("pgx", dbURI)
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cert_user (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL
);`)

	defer dbContainer.Terminate(context.Background())

	repo := &storage.Storage{
		UserStorage: postgre.NewUserStorage(db),
	}

	type fields struct {
		service *service.Service
		l       *zap.Logger
	}

	field := fields{
		service: &service.Service{User: service.NewUserService(repo, service.NewValidateService(), zap.NewExample())},
		l:       zap.NewNop(),
	}

	tests := []struct {
		name         string
		fields       fields
		body         string
		wantResponse interface{}
	}{
		{"success", field, `{"first_name": "cert", "last_name": "test"}`, SuccessUserResponse{
			UserID:  1,
			Message: "user created",
		}},
		{"failed", field, `{"first_name": "cer t", "last_name": "test"}`, ErrorResponse{
			Code:    400,
			Message: http.StatusText(http.StatusBadRequest),
		}},
		{"failed", field, `{"first_nam": "cer t", "last_name": "test"}`, ErrorResponse{
			Code:    400,
			Message: http.StatusText(http.StatusBadRequest),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &Handler{
				service: tt.fields.service,
				l:       tt.fields.l,
			}

			req, err := http.NewRequest("POST", "/rest/user", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			handler.CreateUser(recorder, req)

			if recorder.Code == http.StatusOK {
				var response SuccessUserResponse
				if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to parse response body: %v", err)
				}

				if !reflect.DeepEqual(response, tt.wantResponse) {
					t.Errorf("Expected output %v, got %v", tt.wantResponse, response)
				}
			} else {
				var response ErrorResponse
				if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to parse response body: %v", err)
				}

				if !reflect.DeepEqual(response, tt.wantResponse) {
					t.Errorf("Expected output %v, got %v", tt.wantResponse, response)
				}
			}
		})
	}
}
