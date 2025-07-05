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
package category

import (
	"database/sql"
	"fmt"

	"github.com/cgoesche/willdo/internal/database"
)

const (
	categoryTableName        string = "category"
	DefaultCategoryTableName string = "My List"
	categoryTableSchema      string = `CREATE TABLE IF NOT EXISTS category (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	CONSTRAINT name_constraint UNIQUE (name)
);`
)

type Repository struct {
	db        database.IDatabase
	tableName string
}

func NewRepositoryService(db database.IDatabase) *Repository {
	return &Repository{
		tableName: categoryTableName,
		db:        db,
	}
}

func (r *Repository) Init() error {
	// Create the 'category' table
	_, err := r.db.RawQuery(categoryTableSchema)
	if err != nil {
		return err
	}

	// Create a default category row
	var q = fmt.Sprintf("INSERT OR IGNORE INTO %s (id, name, description) VALUES (1, ?, ?);", r.tableName)
	_, err = r.db.RawQuery(q, DefaultCategoryTableName, "This is the default category")
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Create(c Category) (int64, error) {
	var id int64
	var q = fmt.Sprintf("INSERT INTO %s (name, description) VALUES (?, ?)", r.tableName)
	res, err := r.db.RawQuery(q, c.Name, c.Description)
	if err != nil {
		return -1, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *Repository) Delete(id int64) (int64, error) {
	var q = fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.tableName)
	res, err := r.db.RawQuery(q, id)
	if err != nil {
		return -1, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *Repository) DeleteAll() error {
	var q = fmt.Sprintf("DELETE FROM %s", r.tableName)
	_, err := r.db.RawQuery(q)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetAll() (Categories, error) {
	var cats Categories
	var q = fmt.Sprintf("SELECT * FROM %s;", r.tableName)
	rows, err := r.db.RawRowsQuery(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no categories available")
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (r *Repository) GetById(id int64) (Category, error) {
	var c Category
	var q = fmt.Sprintf("SELECT * FROM %s WHERE id = ?;", r.tableName)
	row := r.db.RawRowQuery(q, id)

	err := row.Scan(&c.ID, &c.Name, &c.Description)
	if err == sql.ErrNoRows {
		return c, fmt.Errorf("no category with ID %d found", id)
	}
	return c, nil
}
func (r *Repository) GetByName(name string) (Category, error) {
	var c Category
	var q = fmt.Sprintf("SELECT * FROM %s WHERE name = ?;", r.tableName)
	row := r.db.RawRowQuery(q, name)

	err := row.Scan(&c.ID, &c.Name, &c.Description)
	if err == sql.ErrNoRows {
		return c, fmt.Errorf("no category with Name %s found", name)
	}
	return c, nil
}

func (r *Repository) Update(c Category) (int64, error) {
	var q = fmt.Sprintf("UPDATE %s SET name = ?, description = ? WHERE id = ?;", r.tableName)
	res, err := r.db.RawQuery(q, c.Name, c.Description, c.ID)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
