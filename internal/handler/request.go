package handler

import "github.com/dgrijalva/jwt-go"

// RecipesRequest object to map incoming request for Recipes handler
type RecipesRequest struct {
	Page        uint64   `schema:"page"`
	Term        string   `schema:"term"`
	Ingredients []string `schema:"ingredient"`
}

// CreateRecipeRequest object to map incoming request for Create handler
type RecipeCreateRequest struct {
	Title       string   `validate:"required,min=2" schema:"title"`
	URL         string   `validate:"required,min=10" schema:"url"`
	Thumbnail   string   `schema:"thumbnail"`
	Ingredients []string `validate:"required,max=30,min=1" schema:"ingredients"`
}

// SignUpRequest object to map sign up incoming request
type SignInRequest struct {
	Username string `json:"username" validate:"required,min=1,max=20"`
	Password string `json:"password" validate:"required,min=1,max=32"`
}

// SignUpRequest object to map sign up incoming request
type SignUpRequest struct {
	Email          string `json:"email" validate:"required,email"`
	FullName       string `json:"fullName"`
	Username       string `json:"username" validate:"required,min=5,max=20"`
	Password       string `json:"password" validate:"required,min=8,max=32"`
	RepeatPassword string `json:"repeatPassword" validate:"eqfield=Password"`
}

// Token object to map incoming authorization bearer token
type Token struct {
	UserID   int64  `json:"uid"`
	Username string `json:"uname"`
	jwt.StandardClaims
}
