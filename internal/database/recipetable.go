package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const recipeColumns = "r.id, r.title, r.thumbnail, r.url, r.created_at, r.updated_at"

// RecipeFilters object
type RecipeFilters struct {
	Term        string
	Ingredients []string
}

// RecipeTable object
type RecipeTable struct {
	db       *sql.DB
	name     string
	pageSize uint64
}

// NewRecipeTable object
func NewRecipeTable(db *sql.DB) *RecipeTable {
	return &RecipeTable{
		db:       db,
		name:     "recipe r",
		pageSize: 10,
	}
}

// Get a recipe by id
func (rt *RecipeTable) Get(id uint64) (*Recipe, error) {
	// nolint:gosec
	query := fmt.Sprintf(`SELECT %s FROM recipe r WHERE id = ?`, recipeColumns)

	var rcp Recipe
	if err := rt.db.QueryRow(query, id).Scan(
		&rcp.ID, &rcp.Title, &rcp.Thumbnail, &rcp.URL, &rcp.CreatedAt, &rcp.UpdatedAt,
	); err != nil {
		return nil, err
	}

	ri, err := rt.withIngredients(rcp)
	if err != nil {
		return nil, err
	}

	return &ri[0], nil
}

// Paginate get paginated recipes
func (rt *RecipeTable) Paginate(page uint64, filters *RecipeFilters) (Recipes, int64, error) {
	var args []interface{}
	query := `SELECT DISTINCT r.id, r.title, r.thumbnail, r.url, r.created_at, r.updated_at 
FROM recipe r 
JOIN ingredient i on r.id = i.recipe_id 
WHERE 1=1`

	if filters != nil && filters.Term != "" {
		query += " AND r.title like ?"
		args = append(args, "%"+filters.Term+"%")
	}

	if filters != nil && len(filters.Ingredients) > 0 {
		query += fmt.Sprintf(" AND i.name in (%s)",
			strings.TrimSuffix(strings.Repeat("?,", len(filters.Ingredients)), ","),
		)
		for i := range filters.Ingredients {
			args = append(args, filters.Ingredients[i])
		}
	}
	query += " GROUP BY r.id"

	// count all results before applying limits
	total, err := rt.countGroup(query, args)
	if err != nil {
		return nil, 0, err
	}

	if page > 0 {
		page--
	}
	query += ` LIMIT ?, ?`
	args = append(args, rt.pageSize*page, rt.pageSize)

	rows, err := rt.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}

	var recipes Recipes
	for rows.Next() {
		r := Recipe{}
		if err := rows.Scan(
			&r.ID, &r.Title, &r.Thumbnail, &r.URL, &r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		recipes = append(recipes, r)
	}
	if rows.Err() != nil {
		return recipes, 0, err
	}

	recipes, err = rt.withIngredients(recipes...)
	if err != nil {
		return nil, 0, err
	}

	return recipes, total, nil
}

// Insert a new recipe, returns inserted recipe id
func (rt *RecipeTable) Insert(recipe Recipe) (int64, error) {
	rq := `INSERT INTO recipe (title, thumbnail, url) VALUES (?, ?, ?)`
	// nolint:gosec
	iq := fmt.Sprintf(`INSERT INTO ingredient (recipe_id, name) VALUES %s`,
		strings.TrimSuffix(strings.Repeat("(?, ?),", len(recipe.Ingredients)), ","),
	)

	tx, err := rt.db.Begin()
	if err != nil {
		return 0, err
	}

	// transactions
	var rid int64
	err = func() error {
		// Insert recipe
		res, err := tx.Exec(rq, recipe.Title, recipe.Thumbnail, recipe.URL)
		if err != nil {
			return fmt.Errorf("recipe error, %w", err)
		}

		// Get last inserted id
		rid, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("recipe error, %w", err)
		}

		// Insert recipe ingredients
		var ingredients []interface{}
		for i := range recipe.Ingredients {
			ingredients = append(ingredients, rid, recipe.Ingredients[i].Name)
		}
		_, err = tx.Exec(iq, ingredients...)
		if err != nil {
			return fmt.Errorf("ingredient error, %w", err)
		}

		return nil
	}()

	// Check if any transaction failed to rollback
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	// Commit transactions
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return rid, nil
}

// Get recipe ingredients
func (rt *RecipeTable) withIngredients(recipes ...Recipe) (Recipes, error) {
	if len(recipes) == 0 {
		return recipes, nil
	}

	var args []interface{}
	// nolint:gosec
	query := fmt.Sprintf(`select id, recipe_id, name, created_at, updated_at 
FROM ingredient 
WHERE recipe_id IN (%s)`,
		strings.TrimSuffix(strings.Repeat("?,", len(recipes)), ","),
	)
	for i := range recipes {
		args = append(args, recipes[i].ID)
	}

	rows, err := rt.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		ing := Ingredient{}
		if err := rows.Scan(&ing.ID, &ing.RecipeID, &ing.Name, &ing.CreatedAt, &ing.UpdatedAt); err != nil {
			return nil, err
		}

		for i := range recipes {
			if recipes[i].ID == ing.RecipeID {
				recipes[i].Ingredients = append(recipes[i].Ingredients, ing)
			}
		}
	}
	if rows.Err() != nil {
		return nil, err
	}

	return recipes, nil
}

func (rt *RecipeTable) countGroup(q string, qArgs []interface{}) (int64, error) {
	q = strings.ReplaceAll(q, "\n", " ")
	q = strings.ReplaceAll(q, "\t", " ")
	neqQ := strings.SplitAfter(strings.ToLower(q), " from ")
	if len(neqQ) <= 1 {
		return 0, fmt.Errorf("unable to count, query should have a from statement, %+v", neqQ)
	}
	// nolint:gosec
	cntQ := fmt.Sprintf("SELECT count(*) as count FROM (%s) as countable", q)

	var total int64
	if err := rt.db.QueryRow(cntQ, qArgs...).Scan(&total); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("unable to count, %w", err)
	}

	return total, nil
}
