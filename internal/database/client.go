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
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cgoesche/willdo/app"
	"github.com/cgoesche/willdo/internal/models"
	_ "modernc.org/sqlite"
)

var (
	ErrInitDatabase = errors.New("could not initialize database")
	ErrOpenDatabase = errors.New("could not open database")
	ErrRowNotFound  = errors.New("the query returned no row")
)

const (
	TaskTableName            string = "task"
	CategoryTableName        string = "category"
	DefaultCategoryTableName string = "My List"
	DatabaseFileName         string = app.Name + ".db"
	databaseDriver           string = "sqlite"

	categoryTableSchema string = `CREATE TABLE IF NOT EXISTS category (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	CONSTRAINT name_constraint UNIQUE (name)
);`
	taskTableSchema string = `CREATE TABLE IF NOT EXISTS task (
	id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	status INTEGER NOT NULL,
	priority INTEGER NOT NULL,
	category INTEGER NOT NULL,
	isfavorite INTEGER NOT NULL,
	FOREIGN KEY (category) REFERENCES category(id) ON DELETE CASCADE
);`
)

type Client struct {
	db *sql.DB
}

func NewClient() *Client {
	return &Client{
		db: &sql.DB{},
	}
}

func (c *Client) InitDB(path string) (err error) {
	path = fmt.Sprintf("file:%s?_foreign_keys=on", path)

	c.db, err = sql.Open(databaseDriver, path)
	if err != nil {
		return err
	}

	// Create the 'category' table
	_, err = c.db.ExecContext(context.Background(), categoryTableSchema)
	if err != nil {
		return err
	}

	// Create a default category table
	_, err = c.db.Exec("INSERT OR IGNORE INTO category (id, name, description) VALUES (1, ?, ?);",
		DefaultCategoryTableName, "This is the default category")
	if err != nil {
		return err
	}

	// Create the 'task' table
	_, err = c.db.ExecContext(context.Background(), taskTableSchema)
	if err != nil {
		return err
	}
	return nil
}

// Inserts a row in the appropriate table where the latter is inferred from the
// type of the passed row
func (c *Client) InsertRow(row any) (int64, error) {
	switch v := row.(type) {
	case models.Category:
		return c.addCategory(v)
	case models.Task:
		return c.addTask(v)
	default:
		return -1, fmt.Errorf("invalid type for provided row")
	}
}

// Delete a row from the specified table
func (c *Client) DeleteRowFromID(tableName string, id int64) (int64, error) {
	q := fmt.Sprintf("DELETE FROM %s WHERE id = %d", tableName, id)
	res, err := c.db.Exec(q)
	if err != nil {
		return -1, err
	}

	id, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return id, nil
}
