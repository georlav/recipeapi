package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Recipe object
type RecipeTable struct {
	pageSize    uint64
	defaultCols []string
	db          *sql.DB
}

// NewRecipe creates a recipe object
func NewRecipeTable(sqlDB *sql.DB) *RecipeTable {
	return &RecipeTable{
		defaultCols: []string{"r.id", "r.title", "r.thumbnail", "r.url", "r.created_at", "r.updated_at"},
		pageSize:    10,
		db:          sqlDB,
	}
}

// Get a recipe by id
func (r *RecipeTable) Get(id string) (*Recipe, error) {
	// nolint:gosec
	query := fmt.Sprintf(
		`SELECT %s FROM recipe r WHERE id = ?`,
		strings.Join(r.defaultCols, ","),
	)

	var rcp Recipe
	if err := r.db.QueryRow(query, id).Scan(
		&rcp.ID, &rcp.Title, &rcp.Thumbnail, &rcp.URL, &rcp.CreatedAt, &rcp.UpdatedAt,
	); err != nil {
		return nil, err
	}

	ri, err := r.withIngredients(rcp)
	if err != nil {
		return nil, err
	}

	return &ri[0], nil
}

// Paginate get paginated recipes
func (r *RecipeTable) Paginate(page uint64, flt *Filters) (Recipes, int64, error) {
	var args []interface{}
	query := `SELECT DISTINCT r.id, r.title, r.thumbnail, r.url, r.created_at, r.updated_at 
FROM recipe r 
JOIN ingredient i on r.id = i.recipe_id 
WHERE 1=1`

	if flt != nil && flt.Term != "" {
		query += " AND r.title like ?"
		args = append(args, "%"+flt.Term+"%")
	}

	if flt != nil && len(flt.Ingredients) > 0 {
		query += fmt.Sprintf(" AND i.name in (%s)",
			strings.TrimSuffix(strings.Repeat("?,", len(flt.Ingredients)), ","),
		)
		for i := range flt.Ingredients {
			args = append(args, flt.Ingredients[i])
		}
	}
	query += " GROUP BY r.id"

	// count all results before applying limits
	total, err := r.countGroup(query, args)
	if err != nil {
		return nil, 0, err
	}

	if page > 0 {
		page--
	}
	query += ` LIMIT ?, ?`
	args = append(args, r.pageSize*page, r.pageSize)

	rows, err := r.db.Query(query, args...)
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

	recipes, err = r.withIngredients(recipes...)
	if err != nil {
		return nil, 0, err
	}

	return recipes, total, nil
}

// Insert updates or insert a new recipe
func (r *RecipeTable) Insert(recipe Recipe) error {
	rq := `INSERT INTO recipe (title, thumbnail, url) VALUES (?, ?, ?)`
	// nolint:gosec
	iq := fmt.Sprintf(`INSERT INTO ingredient (recipe_id, name) VALUES %s ON DUPLICATE KEY UPDATE name = name`,
		strings.TrimSuffix(strings.Repeat("(?, ?),", len(recipe.Ingredients)), ","),
	)

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// transactions
	err = func() error {
		// Insert recipe
		res, err := tx.Exec(rq, recipe.Title, recipe.Thumbnail, recipe.URL)
		if err != nil {
			return fmt.Errorf("recipe error, %w", err)
		}

		// Get last inserted id
		rid, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("recipe error, %w", err)
		}

		// Insert recipe ingredients
		var ingredients []interface{}
		for i := range recipe.Ingredients {
			ingredients = append(ingredients, rid, recipe.Ingredients[i])
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
			return err
		}
		return err
	}

	// Commit transactions
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Get recipe ingredients
func (r *RecipeTable) withIngredients(recipes ...Recipe) (Recipes, error) {
	if len(recipes) == 0 {
		return recipes, nil
	}

	var args []interface{}
	// nolint:gosec
	query := fmt.Sprintf(`select id, recipe_id, name FROM ingredient WHERE recipe_id IN (%s)`,
		strings.TrimSuffix(strings.Repeat("?,", len(recipes)), ","),
	)
	for i := range recipes {
		args = append(args, recipes[i].ID)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			ID         int64
			recipeID   int64
			ingredient string
		)

		if err := rows.Scan(&ID, &recipeID, &ingredient); err != nil {
			return nil, err
		}

		for i := range recipes {
			if recipes[i].ID == fmt.Sprintf(`%d`, recipeID) {
				recipes[i].Ingredients = append(recipes[i].Ingredients, ingredient)
			}
		}
	}
	if rows.Err() != nil {
		return nil, err
	}

	return recipes, nil
}

func (r *RecipeTable) countGroup(q string, qArgs []interface{}) (int64, error) {
	q = strings.ReplaceAll(q, "\n", " ")
	q = strings.ReplaceAll(q, "\t", " ")
	neqQ := strings.SplitAfter(strings.ToLower(q), " from ")
	if len(neqQ) <= 1 {
		return 0, fmt.Errorf("unable to count, query should have a from statement, %+v", neqQ)
	}
	// nolint:gosec
	cntQ := fmt.Sprintf("SELECT count(*) as count FROM (%s) as countable", q)

	var total int64
	if err := r.db.QueryRow(cntQ, qArgs...).Scan(&total); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("unable to count, %w", err)
	}

	return total, nil
}
