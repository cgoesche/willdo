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
package task

import (
	"database/sql"
	"fmt"

	"github.com/cgoesche/willdo/internal/database"
)

const (
	taskTableName   string = "task"
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

type Repository struct {
	db        database.IDatabase
	tableName string
}

func NewRepositoryService(db database.IDatabase) *Repository {
	return &Repository{
		tableName: taskTableName,
		db:        db,
	}
}

func (r *Repository) Init() error {
	// Create the 'task' table
	res, err := r.db.RawQuery(taskTableSchema)
	if err != nil {
		return err
	}

	res.LastInsertId()
	return nil
}

func (r *Repository) Create(t Task) (int64, error) {
	var id int64
	var q = fmt.Sprintf("INSERT INTO %s (title, description, status, priority, category, isfavorite) VALUES (?, ?, ?, ?, ?, ?)", r.tableName)
	res, err := r.db.RawQuery(q,
		t.Title, t.Description, t.Status, t.Priority, t.Category, t.IsFavorite)
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

func (r *Repository) DeleteAllByCategory(cat int64) (int64, error) {
	var tasks Tasks
	var id int64
	// Retrieve all tasks from the specified category
	tasks, err := r.GetAllByCategory(cat)
	if err != nil {
		return -1, err
	}

	for i := 0; i <= len(tasks)-1; i++ {
		if id, err = r.Delete(tasks[i].ID); err != nil {
			return tasks[i].ID, err
		}
	}
	return id, nil
}

func (r *Repository) GetAll() (Tasks, error) {
	var tasks Tasks
	var q = fmt.Sprintf("SELECT * FROM %s;", r.tableName)
	rows, err := r.db.RawRowsQuery(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t Task

		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Category, &t.IsFavorite)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tasks available")
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *Repository) GetById(id int64) (Task, error) {
	var t Task
	var q = fmt.Sprintf("SELECT * FROM %s WHERE id = ?;", r.tableName)
	row := r.db.RawRowQuery(q, id)

	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Category, &t.IsFavorite)
	if err == sql.ErrNoRows {
		return t, fmt.Errorf("no task with ID %d found", id)
	}
	return t, nil
}

func (r *Repository) GetAllByCategory(cat int64) (Tasks, error) {
	var tasks Tasks
	var q = fmt.Sprintf("SELECT * FROM %s WHERE category = ?", r.tableName)
	rows, err := r.db.RawRowsQuery(q, cat)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s table: %v", r.tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Category, &t.IsFavorite)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *Repository) Update(t Task) (int64, error) {
	var q = fmt.Sprintf("UPDATE %s SET title = ?, description = ?, status = ?, priority = ?, category = ?, isfavorite = ? WHERE id = ?;", r.tableName)
	res, err := r.db.RawQuery(q, t.Title, t.Description, t.Status,
		t.Priority, t.Category, t.IsFavorite, t.ID)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
