package recipe

// Repository you need to impl the following methods to create
// a recipe repo of any kind of provider
type Repository interface {
	GetOne(id string) (Recipe, error)
	GetMany(qp QueryParams) (Recipes, error)
}
