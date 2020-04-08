package handler

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/georlav/recipeapi/internal/database"
)

// RecipeResponse recipe response object
type RecipesResponse struct {
	Data     *RecipeResponseItems `json:"data"`
	Metadata Metadata             `json:"metadata"`
}

// Metadata
type Metadata struct {
	Total int64
}

// RecipeResponseItem object to map recipe items
type RecipeResponseItems []RecipeResponseItem

// RecipeResponseItem object to map a recipe item
type RecipeResponseItem struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Href        string             `json:"href"`
	Ingredients IngredientResponse `json:"ingredients"`
	Thumbnail   string             `json:"thumbnail"`
	CreatedAt   string             `json:"createdAt"`
	UpdatedAt   string             `json:"updatedAt"`
}

// IngredientResponseItem object to map single ingredient
type IngredientResponseItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// IngredientResponseItem object to map slice of ingredients
type IngredientResponse []IngredientResponseItem

// UserProfileResponse object to map user profile response
type UserProfileResponse struct {
	ID        int64
	Username  string `json:"username"`
	FullName  string `json:"fullName"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// NewUserProfileResponse creates a new UserProfileResponse object
func NewUserProfileResponse(u database.User) UserProfileResponse {
	return UserProfileResponse{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		Email:     u.Email,
		Active:    u.Active,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// TokenResponse map token response
type TokenResponse struct {
	Token string `json:"token"`
}

// ErrorResponse object to map error response
type ErrorResponse struct {
	Message       string `json:"error"`
	StatusCode    int    `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

// EncodeEntity generic function to map a database entity to a response object
// All fields in entity and response that needed to be mapped should be public
// Response MUST be a pointer to a response struct
func EncodeEntity(entity, response interface{}) error {
	var b bytes.Buffer
	if err := gob.NewEncoder(&b).Encode(entity); err != nil {
		return err
	}

	return gob.NewDecoder(bytes.NewBuffer(b.Bytes())).Decode(response)
}

// EncodeEntities generic function to map a slice of database entities to a response field of a response struct
// All entity fields that will be mapped need to be public, response target field should be same type with entities
// slice, response MUST be a pointer to a response struct
func EncodeEntities(entities, response interface{}, responseField string) error {
	eVal := reflect.ValueOf(entities)
	if eVal.Kind() == reflect.Ptr {
		eVal = eVal.Elem()
	}
	if eVal.Kind() != reflect.Slice {
		return fmt.Errorf("entities argument is expected to be a slice got %s", eVal.Kind())
	}

	respVal := reflect.ValueOf(response)
	if respVal.Kind() != reflect.Ptr {
		return fmt.Errorf("response argument should be a pointer to a response")
	}
	if respVal.Kind() == reflect.Ptr {
		respVal = respVal.Elem()
	}

	respField, ok := respVal.Type().FieldByName(responseField)
	if !ok {
		return fmt.Errorf("%s field could not be found in response. %v", responseField, respField)
	}

	tmpResp := map[string]interface{}{
		responseField: entities,
	}
	b, err := json.Marshal(tmpResp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, response); err != nil {
		return err
	}

	return nil
}
