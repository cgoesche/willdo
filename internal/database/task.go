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
	"database/sql"
	"fmt"

	"github.com/cgoesche/willdo/internal/models"
)

func (c *Client) addTask(t models.Task) (int64, error) {
	res, err := c.db.Exec("INSERT INTO task (title, description, status, priority, category, isfavorite) VALUES (?, ?, ?, ?, ?, ?)",
		t.Title, t.Description, t.Status, t.Priority, t.Category, t.IsFavorite)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (c *Client) DeleteAllTasks() error {
	q := fmt.Sprintf("DELETE FROM %s", TaskTableName)
	_, err := c.db.Query(q)
	if err != nil {
		return fmt.Errorf("failed to delete all tasks")
	}
	return nil
}

func (c *Client) DeleteTask(id int64) error {
	_, err := c.DeleteRowFromID(TaskTableName, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteTasksFromCategory(cat int64) error {
	var tasks models.Tasks
	// Retrieve all tasks from the specified category
	tasks, err := c.QueryTasksFromCategory(cat)
	if err != nil {
		return err
	}

	for i := 0; i <= len(tasks)-1; i++ {
		if err := c.DeleteTask(tasks[i].ID); err != nil {
			return fmt.Errorf("failed to delete task from category")
		}
	}
	return nil
}

// Queries
func (c *Client) QueryAllTasks() (models.Tasks, error) {
	var tasks models.Tasks

	q := fmt.Sprintf("SELECT * FROM %s;", TaskTableName)
	rows, err := c.db.Query(q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t models.Task

		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Category, &t.IsFavorite)
		if err == sql.ErrNoRows {
			return nil, ErrRowNotFound
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (c *Client) QueryTasksFromCategory(cat int64) (models.Tasks, error) {
	var tasks models.Tasks

	q := fmt.Sprintf("SELECT * FROM %s WHERE category = %d", TaskTableName, cat)
	rows, err := c.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("failed to query %s table: %v", TaskTableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Category, &t.IsFavorite)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (c *Client) QueryTaskFromID(id int64) (models.Task, error) {
	var t models.Task

	q := fmt.Sprintf("SELECT * FROM %s WHERE id = %d;", TaskTableName, id)
	row := c.db.QueryRow(q)

	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Category, &t.IsFavorite)
	if err == sql.ErrNoRows {
		return t, ErrRowNotFound
	}
	return t, nil
}

func (c *Client) updateTaskStatus(i int, s models.Status) (int64, error) {
	q := fmt.Sprintf("UPDATE %s SET status = %d WHERE id = %d;", TaskTableName, s, i)
	res, err := c.db.Exec(q)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (c *Client) updateTaskFavoriteStatus(i int64, f int) (int64, error) {
	q := fmt.Sprintf("UPDATE %s SET isfavorite = %d WHERE id = %d;", TaskTableName, f, i)
	res, err := c.db.Exec(q)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (c *Client) CompleteTask(i int) (int64, error) {
	return c.updateTaskStatus(i, models.Done)
}

func (c *Client) StartTask(i int) (int64, error) {
	return c.updateTaskStatus(i, models.Doing)
}

func (c *Client) ResetTask(i int) (int64, error) {
	return c.updateTaskStatus(i, models.ToDo)
}

func (c *Client) ToggleTaskFavoriteStatus(i int64, f int) (int64, error) {
	return c.updateTaskFavoriteStatus(i, f)
}
