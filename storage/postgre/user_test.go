package postgre_test

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zhayt/cert-tz/model"
	"github.com/zhayt/cert-tz/storage/postgre"
)

func TestUserStorage_CreateUser(t *testing.T) {
	// Create a new mock database and obtain a mock instance and a SQL mock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new UserStorage instance with the mock database
	storage := postgre.NewUserStorage(sqlx.NewDb(db, "sqlmock"))

	// Define the expected query and result
	expectedQuery := `INSERT INTO cert_user (first_name, last_name) VALUES ($1, $2) RETURNING id`
	expectedUserID := uint64(1)

	// Set up the mock to expect the query and return the result
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs("Cert", "User").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedUserID))

	// Call the CreateUser method
	user := model.User{
		FirstName: "Cert",
		LastName:  "User",
	}
	userID, err := storage.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)

	// Ensure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserStorage_GetUser(t *testing.T) {
	// Create a new mock database and obtain a mock instance and a SQL mock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new UserStorage instance with the mock database
	storage := postgre.NewUserStorage(sqlx.NewDb(db, "sqlmock"))

	// Define the expected query and result
	expectedQuery := `SELECT * FROM cert_user WHERE id = $1`
	expectedUser := model.User{
		ID:        1,
		FirstName: "Cert",
		LastName:  "User",
		// Set other fields as needed
	}

	// Set up the mock to expect the query and return the result
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(expectedUser.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name"}).
			AddRow(expectedUser.ID, expectedUser.FirstName, expectedUser.LastName))

	// Call the GetUser method
	user, err := storage.GetUser(context.Background(), expectedUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	// Ensure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserStorage_UpdateUser(t *testing.T) {
	// Create a new mock database and obtain a mock instance and a SQL mock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new UserStorage instance with the mock database
	storage := postgre.NewUserStorage(sqlx.NewDb(db, "sqlmock"))

	// Define the expected query and result
	expectedQuery := `UPDATE cert_user SET first_name=$1, last_name=$2 WHERE id=$3 RETURNING id`
	expectedUserID := uint64(1)

	// Set up the mock to expect the query and return the result
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs("Cert", "User", expectedUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedUserID))

	// Call the UpdateUser method
	user := model.User{
		ID:        expectedUserID,
		FirstName: "Cert",
		LastName:  "User",
	}
	userID, err := storage.UpdateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)

	// Ensure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserStorage_DeleteUser(t *testing.T) {
	// Create a new mock database and obtain a mock instance and a SQL mock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create a new UserStorage instance with the mock database
	storage := postgre.NewUserStorage(sqlx.NewDb(db, "sqlmock"))

	// Define the expected query
	expectedQuery := `DELETE FROM cert_user WHERE id = $1`

	// Set up the mock to expect the query
	mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the DeleteUser method
	err = storage.DeleteUser(context.Background(), 1)
	assert.NoError(t, err)

	// Ensure that all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserStorage_UpdateUser1(t *testing.T) {
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
	if err != nil {
		t.Error(err)
	}
	_, err = db.Exec(`INSERT INTO cert_user (first_name, last_name) VALUES ('Certtest', 'Cert')`)
	if err != nil {
		t.Error(err)
	}

	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{"success", args{context.Background(), model.User{ID: 1, FirstName: "Updated", LastName: "Cert"}}, 1},
		{"failed", args{context.Background(), model.User{ID: 2, FirstName: "Updated", LastName: "Cert"}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := postgre.NewUserStorage(db)
			got, _ := r.UpdateUser(tt.args.ctx, tt.args.user)
			assert.Equalf(t, tt.want, got, "UpdateUser(%v, %v)", tt.args.ctx, tt.args.user)
		})
	}
}
