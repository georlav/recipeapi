package recipe

// Repository you need to implement the following methods to create a recipe repository
type Repository interface {
	GetOne(id string) (r Recipe, err error)
	GetMany(qp QueryParams) (r Recipes, total int64, err error)
}
