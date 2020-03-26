package handler

import "github.com/dgrijalva/jwt-go"

// RecipesRequest object to map incoming request for Recipes handler
type RecipesRequest struct {
	Page        uint64   `schema:"page" validate:"omitempty,min=1"`
	Term        string   `schema:"term" validate:"omitempty,min=3"`
	Ingredients []string `schema:"ingredient" validate:"omitempty,max=5"`
}

// CreateRecipeRequest object to map incoming request for Create handler
type RecipeCreateRequest struct {
	Title       string   `json:"title" validate:"required,min=2"`
	URL         string   `json:"url" validate:"required,min=10"`
	Thumbnail   string   `json:"thumbnail"`
	Ingredients []string `json:"ingredients" validate:"required,max=30,min=1"`
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
