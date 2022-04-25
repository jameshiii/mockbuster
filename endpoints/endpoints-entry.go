package endpoints

import (
	"database/sql"
)

var sqlDb *sql.DB

func New(db *sql.DB) {
	sqlDb = db
}
