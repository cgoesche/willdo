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
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	lists    []list.Model
	dbClient *database.Client
}

func initialModel() model {
	return model{
		lists: make([]list.Model, 0),
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
		case msg.String() == "s":
			return m.startTask()
		case msg.String() == "r":
			return m.resetTask()
		}

	case tea.WindowSizeMsg:
		//h, v := docStyle.GetFrameSize()
		m.lists[0].SetSize(50, 15)
	}

	var cmd tea.Cmd
	m.lists[0], cmd = m.lists[0].Update(msg)
	return m, cmd
}

func (m model) View() string {
	var s StatusBar
	statusBar := s.RenderStatusBar(m.lists[0].Items())
	content := docStyle.Render(m.lists[0].View() + "\n" + statusBar)

	return content
}

func Run(client *database.Client) {
	m := initialModel()
	m.dbClient = client

	d := newTaskItemDelegate()

	defaultList := list.New([]list.Item{}, d, 0, 0)
	defaultList.SetStatusBarItemName("task", "tasks")
	m.lists = []list.Model{defaultList}

	l, err := getTaskListItems(client)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	m.lists[0].SetItems(l)
	m.lists[0].Title = "Default"
	m.lists[0].Styles = styles.DefaultStyles()
	m.lists[0].StatusMessageLifetime = 3 * time.Second

	p := tea.NewProgram(m)
	_, err = p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) updateListItems() error {
	l, err := getTaskListItems(m.dbClient)
	if err != nil {
		return err
	}
	m.lists[0].SetItems(l)
	return nil
}
