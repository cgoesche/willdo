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
	"fmt"
	"strings"

	"github.com/cgoesche/willdo/internal/database"
)

type Service struct {
	repo *Repository
}

func NewService(db database.IDatabase) *Service {
	r := NewRepositoryService(db)
	return &Service{
		repo: r,
	}
}

// Initialisez the underlying database and makes sure that
// the required 'task' table is created
func (s *Service) InitRepo() error {
	return s.repo.Init()
}

func (s *Service) Create(t Task) (int64, error) {
	var id int64

	if err := validateTask(t); err != nil {
		return -1, err
	}

	id, err := s.repo.Create(t)
	if err != nil {
		return -1, fmt.Errorf("failed to add task, %v", err)
	}
	return id, nil
}

func (s *Service) Delete(id int64) (int64, error) {
	if id <= 0 {
		return -1, fmt.Errorf("invalid task ID (%d), needs to be > 0", id)
	}

	id, err := s.repo.Delete(id)
	if err != nil {
		return -1, fmt.Errorf("failed to delete task %d, %v", id, err)
	}
	return id, nil
}

func (s *Service) DeleteAll() error {
	if err := s.repo.DeleteAll(); err != nil {
		return fmt.Errorf("failed to delete all tasks, %v", err)
	}
	return nil
}

func (s *Service) DeleteAllByCategory(cat int64) (int64, error) {
	if cat <= 0 {
		return -1, fmt.Errorf("invalid category ID (%d), needs to be > 0", cat)
	}

	id, err := s.repo.DeleteAllByCategory(cat)
	if err != nil {
		return -1, fmt.Errorf("failed to delete all tasks")
	}
	return id, nil
}

func (s *Service) GetAll() (Tasks, error) {
	var tasks Tasks
	tasks, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) GetById(id int64) (Task, error) {
	var task Task

	if id <= 0 {
		return task, fmt.Errorf("invalid task ID (%d), needs to be > 0", id)
	}

	task, err := s.repo.GetById(id)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (s *Service) GetAllByCategory(cat int64) (Tasks, error) {
	if cat <= 0 {
		return nil, fmt.Errorf("invalid category ID (%d), needs to be > 0", cat)
	}

	var tasks Tasks
	tasks, err := s.repo.GetAllByCategory(cat)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *Service) Update(t Task) (int64, error) {
	if err := validateTask(t); err != nil {
		return -1, err
	}

	id, err := s.repo.Update(t)
	if err != nil {
		return id, err
	}
	return id, nil
}

func validateTask(t Task) error {
	if len(strings.TrimSpace(t.Title)) == 0 {
		return fmt.Errorf("title cannot be empty")
	} else if t.Category <= 0 {
		return fmt.Errorf("invalid category ID %d", t.Category)
	}
	return nil
}
