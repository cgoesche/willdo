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
package bubbletea

import (
	"fmt"
	"strconv"

	"github.com/cgoesche/willdo/internal/bubbletea/styles"
	"github.com/cgoesche/willdo/internal/models"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type StatsBar struct {
	Todo     int
	Doing    int
	Done     int
	Progress float64
}

func (s *StatsBar) UpdateStatsBar() error {
	return nil
}

func (s *StatsBar) storeTaskStats(l []list.Item) error {
	for _, v := range l {
		i, ok := v.(taskListItem)
		if !ok {
			return fmt.Errorf("Error")
		}

		switch i.Status() {
		case int64(models.ToDo):
			s.Todo++
		case int64(models.Doing):
			s.Doing++
		case int64(models.Done):
			s.Done++

		}
	}

	if len(l) == 0 {
		s.Progress = 0
	} else {
		s.Progress = (float64(s.Done) / float64(len(l))) * 100
	}
	return nil
}

func (s *StatsBar) RenderStatsBar(l []list.Item) string {
	s.storeTaskStats(l)

	todo := lipgloss.NewStyle().Foreground(styles.Highlight).Render(strconv.Itoa(s.Todo))
	doing := lipgloss.NewStyle().Foreground(styles.Notice).Render(strconv.Itoa(s.Doing))
	done := lipgloss.NewStyle().Foreground(styles.Special).Render(strconv.Itoa(s.Done))
	p := lipgloss.NewStyle().Foreground(styles.Notice).Render(strconv.Itoa(int(s.Progress)) + "%")

	return fmt.Sprintf("Todo %s • Doing %s • Done %s • Progress %s", todo, doing, done, p)
}
