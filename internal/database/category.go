/*
Copyright © 2025 Christian Goeschel Ndjomouo <cgoesc2@wgu.edu>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package database

import (
	"fmt"

	"github.com/cgoesche/willdo/internal/models"
)

const (
	CategoryTableName        string = "category"
	DefaultCategoryTableName string = "My List"
	categoryTableSchema      string = `CREATE TABLE IF NOT EXISTS category (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	CONSTRAINT name_constraint UNIQUE (name)
);`
)

func (c *Client) addCategory(cat models.Category) (int64, error) {
	res, err := c.db.Exec("INSERT INTO category (name, description) VALUES (?, ?)", cat.Name, cat.Description)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (c *Client) DeleteCategory(id int64) error {
	if err := c.DeleteTasksFromCategory(id); err != nil {
		return err
	}
	_, err := c.DeleteRowFromID(CategoryTableName, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteAllCategories() error {
	q := fmt.Sprintf("DELETE FROM %s", CategoryTableName)
	_, err := c.db.Query(q)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) QueryCategories() (models.Categories, error) {
	var cats models.Categories

	q := fmt.Sprintf("SELECT * FROM %s;", CategoryTableName)
	rows, err := c.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cat models.Category

		err := rows.Scan(&cat.ID, &cat.Name, &cat.Description)
		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, nil
}
