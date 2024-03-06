package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/zhayt/cert-tz/internal/model"
	"github.com/zhayt/cert-tz/internal/service"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// CreateUser creates a new user.
//
// @Summary Create a user
// @Description Creates a new user.
// @Tags User
// @Accept json
// @Produce json
// @Param input body model.User true "User object"
// @Success 200 {object} SuccessUserResponse "User ID"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeout)
	defer cancel()

	var user model.User

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.l.Error("create user: decode error", zap.Error(err))
		h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	userID, err := h.service.User.CreateUser(ctx, user)
	if err != nil {
		h.l.Error("CreateUser error", zap.Error(err))

		if errors.Is(err, service.ErrInvalidData) {
			h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		h.respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	h.l.Info("user created", zap.Uint64("userID", userID))
	h.respondWithSuccess(w, userID, "user created")
}

// GetUser retrieves a user by ID.
//
// @Summary Get a user
// @Description Retrieves a user by ID.
// @Tags User
// @Accept json
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} model.User "Successful response"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/{id} [get]
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeout)
	defer cancel()

	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || userID <= 0 {
		h.l.Error("get user: url param error", zap.Error(err))
		h.respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	user, err := h.service.User.GetUser(ctx, uint64(userID))
	if err != nil {
		h.l.Error("GetUser error", zap.Error(err))

		if errors.Is(err, sql.ErrNoRows) {
			h.respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

		h.respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	h.l.Info("user found", zap.Uint64("userID", user.ID))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates an existing user.
//
// @Summary Update a user
// @Description Updates an existing user.
// @Tags User
// @Accept json
// @Produce json
// @Param id path integer true "User ID"
// @Param input body model.User true "User object"
// @Success 200 {object} model.User "Successful response"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 404 {object} ErrorResponse "Not found"
// @Router /user/{id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeout)
	defer cancel()

	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || userID <= 0 {
		h.l.Error("update user: url param error", zap.Error(err))

		h.respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	var user model.User

	err = json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()
	if err != nil {
		h.l.Error("create user: decode error", zap.Error(err))

		h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user.ID = uint64(userID)

	_, err = h.service.User.UpdateUser(ctx, user)
	if err != nil {
		h.l.Error("UpdateUser error", zap.Error(err))
		if errors.Is(err, service.ErrInvalidData) {
			h.respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		h.respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	h.l.Info("user updated", zap.Uint64("userID", user.ID))
	h.respondWithSuccess(w, user.ID, "user updated")
}

// DeleteUser deletes a user by ID.
//
// @Summary Delete a user
// @Description Deletes a user by ID.
// @Tags User
// @Accept json
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} SuccessUserResponse "Successful response"
// @Failure 404 {object} ErrorResponse "Not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), _defaultContextTimeout)
	defer cancel()

	userID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || userID <= 0 {
		h.l.Error("update user: url param error", zap.Error(err))

		h.respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	err = h.service.User.DeleteUser(ctx, uint64(userID))
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		h.l.Error("DeleteUser error", zap.Error(err))
	}

	h.l.Info("User deleted", zap.Int("userID", userID))
	h.respondWithSuccess(w, uint64(userID), "User deleted")
}
