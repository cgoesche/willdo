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
