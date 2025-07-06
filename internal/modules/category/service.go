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

func (s *Service) Create(c Category) (int64, error) {
	var id int64

	if err := validateCategory(c); err != nil {
		return -1, err
	}

	id, err := s.repo.Create(c)
	if err != nil {
		return -1, fmt.Errorf("failed to add category, %v", err)
	}
	return id, nil
}

func (s *Service) Delete(id int64) (int64, error) {
	if id <= 0 {
		return -1, fmt.Errorf("invalid category ID (%d), needs to be > 0", id)
	}

	id, err := s.repo.Delete(id)
	if err != nil {
		return -1, fmt.Errorf("failed to delete category %d, %v", id, err)
	}
	return id, nil
}

func (s *Service) DeleteAll() error {
	if err := s.repo.DeleteAll(); err != nil {
		return fmt.Errorf("failed to delete all categories, %v", err)
	}
	return nil
}

func (s *Service) GetAll() (Categories, error) {
	var cats Categories
	cats, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *Service) GetById(id int64) (Category, error) {
	var cat Category

	if id <= 0 {
		return cat, fmt.Errorf("invalid category ID (%d), needs to be > 0", id)
	}

	cat, err := s.repo.GetById(id)
	if err != nil {
		return cat, err
	}
	return cat, nil
}

func (s *Service) Update(c Category) (int64, error) {
	if err := validateCategory(c); err != nil {
		return -1, err
	}

	id, err := s.repo.Update(c)
	if err != nil {
		return id, err
	}
	return id, nil
}

func validateCategory(c Category) error {
	if len(strings.TrimSpace(c.Name)) == 0 {
		return fmt.Errorf("title cannot be empty")
	}
	return nil
}
