package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/georlav/recipeapi/internal/database"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

// User godoc
// @Summary user profile
// @Description Get user profile info
// @ID user-profile
// @Accept  application/x-www-form-urlencoded
// @Produce  json
// @Success 200 {object} handler.UserProfileResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Security ApiKeyAuth
// @Router /user [get]
func (h Handler) User(w http.ResponseWriter, r *http.Request) {
	token, err := h.getToken(r)
	if err != nil {
		h.respondError(w, APIError{Message: err.Error(), StatusCode: http.StatusUnauthorized})
		return
	}

	user, err := h.db.User.Get(uint64(token.UserID))
	if err != nil {
		h.respondError(w, APIError{Message: "unknown user", StatusCode: http.StatusNotFound})
		return
	}

	h.respond(w, NewUserProfileResponse(*user), http.StatusOK)
}

// SignIn godoc
// @Summary user sign in
// @Description user sign in
// @ID user-sign-in
// @Accept  json
// @Produce  json
// @Param credentials body handler.SignInRequest false "credentials payload"
// @Success 200 {object} handler.TokenResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /user/signin [post]
func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	si := SignInRequest{}
	if err := json.NewDecoder(r.Body).Decode(&si); err != nil {
		h.respondError(w, APIError{Message: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest})
		return
	}

	// Validate sign up request
	v := validator.New()
	if err := v.Struct(si); err != nil {
		h.respondError(w, APIError{Message: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest})
		return
	}

	// Get user
	u, err := h.db.User.GetByUsername(si.Username)
	if err != nil {
		h.respondError(w, APIError{
			Message:    "You have entered an invalid username or password",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	// User exists check password
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(si.Password)); err != nil {
		h.respondError(w, APIError{
			Message:    "You have entered an invalid username or password",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	// User password is correct create token
	token, err := h.newToken(u)
	if err != nil {
		h.respondError(w, err)
		return
	}

	// Respond with a valid token
	resp := TokenResponse{Token: *token}
	h.respond(w, resp, http.StatusOK)
}

// SignUp godoc
// @Summary user sign up
// @Description user sign up
// @ID user-sign-up
// @Accept  json
// @Produce  json
// @Param body body handler.SignUpRequest true "sign up payload"
// @Success 200 {object} handler.TokenResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /user/signup [post]
func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Map request to struct
	u := SignUpRequest{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		h.respondError(w, APIError{Message: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest})
		return
	}

	// Validate sign up request
	v := validator.New()
	if err := v.Struct(u); err != nil {
		h.respondError(w, APIError{Message: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest})
		return
	}

	// Check if username is available
	if _, err := h.db.User.GetByUsername(u.Username); err == nil {
		h.respondError(w, APIError{Message: "Username is taken", StatusCode: http.StatusConflict})
		return
	}

	// Create hash from password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		h.respondError(w, APIError{Message: http.StatusText(http.StatusBadRequest), StatusCode: http.StatusBadRequest})
		return
	}

	// Create user
	uID, err := h.db.User.Insert(database.User{
		Username: u.Username,
		Password: string(hash),
		FullName: u.FullName,
		Email:    u.Email,
		Active:   true,
	})
	if err != nil {
		h.respondError(w, errors.New("failed to create user"))
		return
	}

	// User created, generate a token
	token, err := h.newToken(&database.User{ID: uID, Username: u.Username})
	if err != nil {
		h.respondError(w, err)
		return
	}

	// Respond with a token
	resp := TokenResponse{Token: *token}
	h.respond(w, resp, http.StatusOK)
}
