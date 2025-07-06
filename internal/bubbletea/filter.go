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
package bubbletea

import (
	"fmt"

	"github.com/cgoesche/willdo/internal/modules/task"
	"github.com/charmbracelet/bubbles/list"
)

func (m model) Filter(value any) ([]list.Item, error) {
	var l []list.Item
	var filterFn func(item list.Item, value any) (list.Item, error)

	switch v := value.(type) {
	case task.Status:
		filterFn = statusFilter
	case task.Priority:
		filterFn = priorityFilter
	case task.FavoriteFlag:
		filterFn = favoriteFilter
	default:
		return nil, fmt.Errorf("unknown filter type %v", v)
	}

	var err error
	for _, i := range m.cachedItems {
		if i, err = filterFn(i, value); err != nil {
			continue
		}
		l = append(l, i)
	}
	return l, nil
}

func statusFilter(item list.Item, value any) (list.Item, error) {
	i, ok := item.(taskListItem)
	if !ok {
		return nil, fmt.Errorf("invalid item type")
	}

	v, ok := value.(task.Status)
	if !ok {
		return nil, fmt.Errorf("invalid filter value type")
	}

	if i.Status() != int64(v) {
		return nil, fmt.Errorf("non-matching item")
	}
	return item, nil
}

func priorityFilter(item list.Item, value any) (list.Item, error) {
	i, ok := item.(taskListItem)
	if !ok {
		return nil, fmt.Errorf("invalid item type")
	}

	v, ok := value.(task.Priority)
	if !ok {
		return nil, fmt.Errorf("invalid filter value type")
	}

	if i.Priority() != int64(v) {
		return nil, fmt.Errorf("non-matching item")
	}
	return item, nil
}

func favoriteFilter(item list.Item, value any) (list.Item, error) {
	i, ok := item.(taskListItem)
	if !ok {
		return nil, fmt.Errorf("invalid item type")
	}

	v, ok := value.(task.FavoriteFlag)
	if !ok {
		return nil, fmt.Errorf("invalid filter value type")
	}

	if i.IsFavorite() != int(v) {
		return nil, fmt.Errorf("non-matching item")
	}
	return item, nil
}
