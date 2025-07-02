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
	"github.com/cgoesche/willdo/internal/bubbletea/styles"
	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/models"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	lists            []list.Model
	stats            statsModel
	details          detailsModel
	DbClient         *database.Client
	selectedList     int
	Categories       models.Categories
	CatNameToIDMap   models.CategoryNameToIDMap
	CatIDToNameMap   models.CategoryIDToNameMap
	SelectedCategory int64

	ShowDetails  bool
	ShowStats    bool
	ShowAllTasks bool
}

func InitialModel() model {
	d := newDetailsModel()

	return model{
		details:      d,
		lists:        make([]list.Model, 0),
		selectedList: 0,
		ShowDetails:  true,
		ShowStats:    true,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}
	task := m.lists[m.selectedList].SelectedItem()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case msg.String() == "ctrl+d":
			m.ShowDetails = !m.ShowDetails
			return m, m.lists[m.selectedList].NewStatusMessage("Toggled 'Details' section")
		case msg.String() == "ctrl+s":
			m.ShowStats = !m.ShowStats
			return m, m.lists[m.selectedList].NewStatusMessage("Toggled 'Statistics' section")
		case msg.String() == "ctrl+r":
			return m.refreshTaskList()
		case msg.String() == "d":
			return m.deleteTask(task)
		case msg.String() == "c":
			return m.completeTask(task)
		case msg.String() == "f":
			return m.toggleTaskFavStatus(task)
		case msg.String() == "s":
			return m.startTask(task)
		case msg.String() == "r":
			return m.resetTask(task)
		}

	case tea.WindowSizeMsg:
		h, v := styles.DocStyle.GetFrameSize()
		m.lists[m.selectedList].SetSize(msg.Width-v, msg.Height-(2*h))
	}

	m.lists[m.selectedList], cmd = m.lists[m.selectedList].Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var content string

	content = m.lists[m.selectedList].View() + "\n"

	if m.ShowDetails {
		m.details.updateSelectedItem(m.lists[m.selectedList].SelectedItem())
		content += m.details.View()
	}

	if m.ShowStats {
		m.stats.updateList(m.lists[m.selectedList].Items())
		content += m.stats.View()
	}
	return styles.DocStyle.Render(content)
}

func (m *model) updateListItems() (err error) {
	var l []list.Item
	if m.ShowAllTasks {
		l, err = getAllTaskListItems(m.DbClient)
	} else {
		l, err = getTaskListItemsByCategory(m.DbClient, m.SelectedCategory)
	}

	if err != nil {
		return err
	}
	m.lists[m.selectedList].SetItems(l)
	return nil
}

func (m model) refreshTaskList() (tea.Model, tea.Cmd) {
	if err := m.updateListItems(); err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to refresh list refreshed")
	}
	return m, m.lists[m.selectedList].NewStatusMessage("List refreshed")
}

func (m model) deleteTask(item list.Item) (tea.Model, tea.Cmd) {
	task, ok := item.(taskListItem)
	if !ok {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to delete task")
	}

	err := m.DbClient.DeleteTask(task.ID)
	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to delete task")
	}

	if err = m.updateListItems(); err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Task deleted but failed to update task list")
	}
	return m, m.lists[m.selectedList].NewStatusMessage("Task deleted")
}

func (m model) completeTask(item list.Item) (tea.Model, tea.Cmd) {
	task, ok := item.(taskListItem)
	if !ok {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Done'")
	}

	_, err := m.DbClient.CompleteTask(int(task.ID))
	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Done'")
	}

	if err = m.updateListItems(); err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Task marked as 'Done' but failed to update task list")
	}
	return m, m.lists[m.selectedList].NewStatusMessage("Task marked as 'Done'")
}

func (m model) startTask(item list.Item) (tea.Model, tea.Cmd) {
	task, ok := item.(taskListItem)
	if !ok {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Doing'")
	}

	_, err := m.DbClient.StartTask(int(task.ID))
	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Doing'")
	}

	if err = m.updateListItems(); err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Task marked as 'Doing' but failed to update task list")
	}
	return m, m.lists[m.selectedList].NewStatusMessage("Task marked as 'Doing'")
}

func (m model) resetTask(item list.Item) (tea.Model, tea.Cmd) {
	task, ok := item.(taskListItem)
	if !ok {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Todo'")
	}

	_, err := m.DbClient.ResetTask(int(task.ID))
	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Todo'")
	}

	if err = m.updateListItems(); err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Task marked as 'Todo' but failed to update task list")
	}
	return m, m.lists[m.selectedList].NewStatusMessage("Task marked as 'Todo'")
}

func (m model) toggleTaskFavStatus(item list.Item) (tea.Model, tea.Cmd) {
	task, ok := item.(taskListItem)
	if !ok {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Favorite'")
	}

	var favStatus int
	switch task.IsFav {
	case models.IsFavorite:
		favStatus = models.IsNotFavorite
	case models.IsNotFavorite:
		favStatus = models.IsFavorite
	}

	_, err := m.DbClient.ToggleTaskFavoriteStatus(task.ID, favStatus)
	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Favorite'")
	}

	if err = m.updateListItems(); err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Task marked as 'Favorite' but failed to update task list")
	}
	return m, m.lists[m.selectedList].NewStatusMessage("Task marked as 'Favorite'")
}
