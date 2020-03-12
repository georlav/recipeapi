package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func countGroup(db *sql.DB, q string, qArgs []interface{}) (int64, error) {
	q = strings.ReplaceAll(q, "\n", " ")
	q = strings.ReplaceAll(q, "\t", " ")
	neqQ := strings.SplitAfter(strings.ToLower(q), " from ")
	if len(neqQ) <= 1 {
		return 0, fmt.Errorf("unable to count, query should have a from statement, %+v", neqQ)
	}
	// nolint:gosec
	cntQ := fmt.Sprintf("SELECT count(*) as count FROM (%s) as countable", q)

	var total int64
	if err := db.QueryRow(cntQ, qArgs...).Scan(&total); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("unable to count, %w", err)
	}

	return total, nil
}
