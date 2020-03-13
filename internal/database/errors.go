package database

import (
	"database/sql"
	"errors"
)

var ErrDuplicateEntry = errors.New("item already exists")
var ErrNoRows = sql.ErrNoRows
