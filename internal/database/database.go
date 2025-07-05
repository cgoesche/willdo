/*
Copyright © 2025 Christian Goeschel Ndjomouo <cgoesc2@wgu.edu>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package database

import (
	"database/sql"
	"errors"

	"github.com/cgoesche/willdo/internal/config"
	_ "modernc.org/sqlite"
)

type IDatabase interface {
	Connect(conf config.Database) (*sql.DB, error)
	RawQuery(q string, args ...any) (sql.Result, error)
	RawRowQuery(q string, args ...any) *sql.Row
	RawRowsQuery(q string, args ...any) (*sql.Rows, error)
}

var (
	ErrInitDatabase = errors.New("could not initialize database")
	ErrOpenDatabase = errors.New("could not open database")
	ErrRowNotFound  = errors.New("the query returned no row")
)

func New(conf config.Database) IDatabase {
	var db IDatabase

	switch conf.Type {
	case "sqlite":
		db = NewSQLite(conf)
	default:
		db = NewSQLite(conf)
	}
	return db
}
