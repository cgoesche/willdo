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
	"fmt"

	"github.com/cgoesche/willdo/internal/config"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(conf config.Database) SQLite {
	var s SQLite
	s.db, _ = s.Connect(conf)

	return s
}

func (s SQLite) Connect(conf config.Database) (*sql.DB, error) {
	path := fmt.Sprintf("file:%s?_foreign_keys=on", conf.Filepath)

	var err error
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s SQLite) RawRowQuery(q string, args ...any) *sql.Row {
	return s.db.QueryRow(q, args...)
}

func (s SQLite) RawRowsQuery(q string, args ...any) (*sql.Rows, error) {
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s SQLite) RawQuery(q string, args ...any) (sql.Result, error) {
	return s.db.Exec(q, args...)
}
