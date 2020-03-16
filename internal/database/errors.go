package database

import (
	"database/sql"
	"errors"
)

var ErrDuplicateEntry = errors.New("already exists")
var ErrNoRows = sql.ErrNoRows
