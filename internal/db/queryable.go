package db

// Queryable you need to implement the following methods to support a new database with new queries
type Queryable interface {
	Get(id string) (*Recipe, error)
	Insert(Recipe) error
	Paginate(page uint64, flt *Filters) (Recipes, int64, error)
}
