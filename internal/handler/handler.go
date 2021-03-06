package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/georlav/recipeapi/internal/config"
	"github.com/georlav/recipeapi/internal/database"
	"github.com/georlav/recipeapi/internal/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type contextKey string

const CtxKeyToken contextKey = "token"

type Handler struct {
	db       *database.Database
	cfg      *config.Config
	log      *logger.Logger
	schema   *schema.Decoder
	validate *validator.Validate
}

func NewHandler(db *database.Database, c *config.Config, l *logger.Logger) *Handler {
	return &Handler{
		db:       db,
		cfg:      c,
		log:      l,
		schema:   schema.NewDecoder(),
		validate: validator.New(),
	}
}

func (h *Handler) newToken(u *database.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uname": u.Username,
		"uid":   u.ID,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Duration(h.cfg.Token.TTL) * time.Minute).Unix(),
	})

	tokenSigned, err := token.SignedString([]byte(h.cfg.Token.Secret))
	if err != nil {
		return nil, fmt.Errorf("signature error, %w", err)
	}

	return &tokenSigned, nil
}

func (h *Handler) getToken(r *http.Request) (*Token, error) {
	token, ok := r.Context().Value(CtxKeyToken).(Token)
	if !ok {
		return nil, errors.New("token not present in header")
	}

	return &token, nil
}
