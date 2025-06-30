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
	"os"
	"time"

	"github.com/cgoesche/willdo/internal/bubbletea/styles"
	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/models"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	lists            []list.Model
	dbClient         *database.Client
	listIndex        int
	categories       models.Categories
	selectedCategory int64
}

func initialModel() model {
	return model{
		lists:     make([]list.Model, 0),
		listIndex: 0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case msg.String() == "d":
			return m.deleteTask()
		case msg.String() == "c":
			return m.completeTask()
		case msg.String() == "f":
			return m.toggleTaskFavStatus()
		case msg.String() == "s":
			return m.startTask()
		case msg.String() == "r":
			return m.resetTask()
		}

	case tea.WindowSizeMsg:
		_, v := docStyle.GetFrameSize()
		m.lists[m.listIndex].SetSize(msg.Width-v, 15)
	}

	var cmd tea.Cmd
	m.lists[m.listIndex], cmd = m.lists[m.listIndex].Update(msg)
	return m, cmd
}

func (m model) View() string {
	var s StatsBar

	details := m.RenderTaskDetailsSection()
	statsBar := s.RenderStatsBar(m.lists[m.listIndex].Items())
	content := docStyle.Render(m.lists[m.listIndex].View() + "\n" + details + statsBar)

	return content
}

func Run(client *database.Client, categories models.Categories, categoryID int64) {
	m := initialModel()
	m.dbClient = client
	m.categories = categories
	m.selectedCategory = categoryID

	d := newTaskItemDelegate()

	defaultList := list.New([]list.Item{}, d, 0, 0)
	defaultList.SetStatusBarItemName("task", "tasks")
	m.lists = []list.Model{defaultList, defaultList}

	l, err := getTaskListItemsByCategory(client, m.selectedCategory)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	m.lists[m.listIndex].SetItems(l)
	m.lists[m.listIndex].Title = models.GetCategoryName(m.categories, m.selectedCategory)
	m.lists[m.listIndex].Styles = styles.DefaultStyles()
	m.lists[m.listIndex].StatusMessageLifetime = 3 * time.Second
	m.lists[m.listIndex].SetShowFilter(false)
	m.lists[m.listIndex].SetFilteringEnabled(false)

	p := tea.NewProgram(m)
	_, err = p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m *model) updateListItems() error {
	l, err := getTaskListItemsByCategory(m.dbClient, m.selectedCategory)
	if err != nil {
		return err
	}
	m.lists[m.listIndex].SetItems(l)
	return nil
}

func (m *model) nextList() (tea.Model, tea.Cmd) {
	if m.listIndex+1 <= len(m.lists)-1 {
		m.listIndex = m.listIndex + 1
	}

	_, cmd := m.lists[m.listIndex].Update(nil)
	return m, cmd
}
