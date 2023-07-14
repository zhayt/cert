package postgre

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zhayt/cert-tz/model"
)

type UserStorage struct {
	db *sqlx.DB
}

func (r *UserStorage) CreateUser(ctx context.Context, user model.User) (uint64, error) {
	qr := `INSERT INTO cert_user (first_name, last_name) VALUES ($1, $2) RETURNING id`

	var userID uint64
	if err := r.db.GetContext(ctx, &userID, qr, user.FirstName, user.LastName); err != nil {
		return 0, fmt.Errorf("cannot create user: %w", err)
	}

	return userID, nil
}

func (r *UserStorage) GetUser(ctx context.Context, userID uint64) (model.User, error) {
	qr := `SELECT * FROM cert_user WHERE id = $1`

	var user model.User
	if err := r.db.GetContext(ctx, &user, qr, userID); err != nil {
		return model.User{}, fmt.Errorf("cannot get user: %w", err)
	}

	return user, nil
}

func (r *UserStorage) UpdateUser(ctx context.Context, user model.User) (uint64, error) {
	qr := `UPDATE cert_user SET first_name=$1, last_name=$2 WHERE id=$3 RETURNING id`

	var userID uint64

	if err := r.db.GetContext(ctx, &userID, qr, user.FirstName, user.LastName, user.ID); err != nil {
		return 0, fmt.Errorf("cannot update user: %w", err)
	}

	return userID, nil
}

func (r *UserStorage) DeleteUser(ctx context.Context, userID uint64) error {
	qr := `DELETE FROM cert_user WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, qr, userID); err != nil {
		return fmt.Errorf("cannot delete user: %w", err)
	}

	return nil
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}
