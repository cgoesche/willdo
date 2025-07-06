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
	"github.com/cgoesche/willdo/internal/modules/task"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	separator = "•"
)

type statsBar struct {
	Todo        int
	Doing       int
	Done        int
	Progress    float64
	IsFiltering bool
}

type statsModel struct {
	list []list.Item
	bar  statsBar
}

func (s statsModel) Init() tea.Cmd {
	return nil
}

func (s statsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s statsModel) View() string {
	return s.renderStatsBar()
}

func (s *statsModel) updateList(l []list.Item) {
	s.list = l
}

func (s *statsModel) storeTaskStats() error {
	l := s.list
	for _, v := range l {
		i, ok := v.(taskListItem)
		if !ok {
			return fmt.Errorf("Error")
		}

		switch i.Status() {
		case int64(task.ToDo):
			s.bar.Todo++
		case int64(task.Doing):
			s.bar.Doing++
		case int64(task.Done):
			s.bar.Done++

		}
	}

	if len(l) == 0 {
		s.bar.Progress = 0
	} else {
		s.bar.Progress = (float64(s.bar.Done) / float64(len(l))) * 100
	}
	return nil
}

func (s *statsModel) renderStatsBar() string {
	var todo, doing, done, prog, filterIndicator string
	s.storeTaskStats()

	todo = lipgloss.NewStyle().Foreground(styles.SubtleColor).Render("Todo ")
	doing = lipgloss.NewStyle().Foreground(styles.SubtleColor).Render("Doing ")
	done = lipgloss.NewStyle().Foreground(styles.SubtleColor).Render("Done ")
	prog = lipgloss.NewStyle().Foreground(styles.SubtleColor).Render("Progress ")

	todo += styles.TodoStyle.Render(strconv.Itoa(s.bar.Todo))
	doing += styles.DoingStyle.Render(strconv.Itoa(s.bar.Doing))
	done += styles.DoneStyle.Render(strconv.Itoa(s.bar.Done))

	if s.bar.IsFiltering {
		filterIndicator = lipgloss.NewStyle().Foreground(styles.LightSubtleColor).Render("[ Filter active ]")
	}

	p := int(s.bar.Progress)
	if p >= 70 {
		prog += styles.DoneStyle.Render(strconv.Itoa(int(s.bar.Progress)) + "%")
	} else if p > 30 && p < 70 {
		prog += lipgloss.NewStyle().Foreground(styles.NoticeColor).Render(strconv.Itoa(int(s.bar.Progress)) + "%")
	} else {
		prog += lipgloss.NewStyle().Foreground(styles.WarningColor).Render(strconv.Itoa(int(s.bar.Progress)) + "%")
	}

	return fmt.Sprintf("%s %s %s %s %s %s %s %s", todo, separator, doing, separator, done, separator, prog, filterIndicator)
}
