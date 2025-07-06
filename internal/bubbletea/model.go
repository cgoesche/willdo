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

	"github.com/cgoesche/willdo/internal/bubbletea/keys"
	"github.com/cgoesche/willdo/internal/bubbletea/styles"

	"github.com/cgoesche/willdo/internal/modules/category"
	"github.com/cgoesche/willdo/internal/modules/task"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list        list.Model
	cachedItems []list.Item
	stats       statsModel
	details     detailsModel

	CategoryService    *category.Service
	TaskService        *task.Service
	Categories         category.Categories
	Tasks              task.Tasks
	CatNameToIDMap     category.CategoryNameToIDMap
	CatIDToNameMap     category.CategoryIDToNameMap
	SelectedCategoryID int64

	KeyMap      keys.KeyMap
	IsFiltering bool
	FilterValue any

	ShowDetails   bool
	ShowStats     bool
	ShowAllTasks  bool
	ShowFavorites bool
}

func InitialModel() model {
	d := newDetailsModel()

	return model{
		details:     d,
		list:        list.Model{},
		ShowDetails: true,
		ShowStats:   true,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var statusMsg string
	var cmd tea.Cmd
	cmds := []tea.Cmd{}
	item := m.list.SelectedItem()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.ToggleDetails):
			m.ShowDetails = !m.ShowDetails
			return m, m.list.NewStatusMessage("Toggled 'Details' section")
		case key.Matches(msg, m.KeyMap.ToggleStats):
			m.ShowStats = !m.ShowStats
			return m, m.list.NewStatusMessage("Toggled 'Statistics' section")
		case key.Matches(msg, m.KeyMap.RefreshList):
			return m.updateListItems("List refreshed")
		// Task modification key mappings
		case key.Matches(msg, m.KeyMap.DeleteTask):
			return m.deleteTask(item)
		case key.Matches(msg, keys.DefaultKeyMap.ClearCompletedTasks):
			return m.deleteCompletedTasks()
		case key.Matches(msg, m.KeyMap.CompleteTask):
			return m.completeTask(item)
		case key.Matches(msg, m.KeyMap.ToggleFavStatus):
			return m.toggleTaskFavStatus(item)
		case key.Matches(msg, m.KeyMap.StartTask):
			return m.startTask(item)
		case key.Matches(msg, m.KeyMap.ResetTask):
			return m.resetTask(item)
		// Filter key mapping
		case key.Matches(msg, m.KeyMap.ClearFilter):
			m.IsFiltering = false
			m.FilterValue = nil
			return m.updateListItems("Filter cleared")
		case key.Matches(msg, m.KeyMap.FilterToDo):
			m.IsFiltering = true
			m.FilterValue = task.ToDo
			l, err := m.Filter(m.FilterValue)
			statusMsg = "Filtering todo tasks"
			if err != nil {
				statusMsg = "Failed to filter"
			}
			cmds = append(cmds, m.list.NewStatusMessage(statusMsg))
			cmds = append(cmds, m.list.SetItems(l))
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.KeyMap.FilterFav):
			m.IsFiltering = true
			m.FilterValue = task.IsFavorite
			l, err := m.Filter(m.FilterValue)
			statusMsg = "Filtering favorites"
			if err != nil {
				statusMsg = "Failed to filter"
			}
			cmds = append(cmds, m.list.NewStatusMessage(statusMsg))
			cmds = append(cmds, m.list.SetItems(l))
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.KeyMap.FilterDoing):
			m.IsFiltering = true
			m.FilterValue = task.Doing
			l, err := m.Filter(m.FilterValue)
			statusMsg = "Filtering doing tasks"
			if err != nil {
				statusMsg = "Failed to filter"
			}
			cmds = append(cmds, m.list.NewStatusMessage(statusMsg))
			cmds = append(cmds, m.list.SetItems(l))
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.KeyMap.FilterDone):
			m.IsFiltering = true
			m.FilterValue = task.Done
			l, err := m.Filter(m.FilterValue)
			statusMsg = "Filtering done tasks"
			if err != nil {
				statusMsg = "Failed to filter"
			}
			cmds = append(cmds, m.list.NewStatusMessage(statusMsg))
			cmds = append(cmds, m.list.SetItems(l))
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.KeyMap.NextCategory):
			m.IsFiltering = false
			m.SelectedCategoryID = 0
			l, err := m.Filter(m.FilterValue)
			if err != nil {
				return m, m.list.NewStatusMessage(fmt.Sprintf("Failed to filter, %v", err))
			}
			return m, m.list.SetItems(l)
		}

	case tea.WindowSizeMsg:
		h, v := styles.DocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-v, msg.Height-(2*h))
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var content string

	content = m.list.View() + "\n"

	if m.ShowDetails {
		m.details.updateSelectedItem(m.list.SelectedItem())
		content += m.details.View()
	}

	if m.ShowStats {
		m.stats.updateList(m.list.Items())
		content += m.stats.View()
	}
	return styles.DocStyle.Render(content)
}

func (m model) getAllTaskListItems() ([]list.Item, error) {
	var l []list.Item

	tasks, err := m.TaskService.GetAll()
	if err != nil {
		return nil, err
	}

	l = marshalTaskListItems(tasks)
	return l, nil
}

func (m model) getTaskListItemsByCategory(id int64) ([]list.Item, error) {
	var l []list.Item

	tasks, err := m.TaskService.GetAllByCategory(id)
	if err != nil {
		return nil, err
	}

	l = marshalTaskListItems(tasks)
	return l, nil
}

func (m *model) updateListItems(statusMsg string) (tea.Model, tea.Cmd) {
	var l []list.Item
	var err error

	switch {
	case m.ShowAllTasks:
		l, err = m.getAllTaskListItems()
	default:
		l, err = m.getTaskListItemsByCategory(m.SelectedCategoryID)
	}
	if err != nil {
		return m, m.list.NewStatusMessage(fmt.Sprintf("Failed to refresh list, %v", err))
	}
	m.cachedItems = l

	if m.IsFiltering {

		l, _ = m.Filter(m.FilterValue)
	}

	m.list.SetItems(l)
	return m, m.list.NewStatusMessage(statusMsg)
}

func (m model) deleteTask(item list.Item) (tea.Model, tea.Cmd) {
	task, ok := item.(taskListItem)
	if !ok {
		return m, m.list.NewStatusMessage("Failed to delete task")
	}

	_, err := m.TaskService.Delete(task.ID)
	if err != nil {
		return m, m.list.NewStatusMessage("Failed to delete task")
	}
	return m.updateListItems("Task deleted")
}

func (m model) deleteCompletedTasks() (tea.Model, tea.Cmd) {
	var l []list.Item

	l, err := m.Filter(task.Done)
	if err != nil {
		return m, m.list.NewStatusMessage("Failed to clear completed tasks")
	}

	for _, i := range l {
		task, ok := i.(taskListItem)
		if !ok {
			return m, m.list.NewStatusMessage("Failed to clear completed tasks")
		}

		_, err := m.TaskService.Delete(task.ID)
		if err != nil {
			return m, m.list.NewStatusMessage("Failed to clear completed tasks")
		}
	}
	return m.updateListItems("Cleared completed tasks")
}

func (m model) completeTask(item list.Item) (tea.Model, tea.Cmd) {
	t, ok := item.(taskListItem)
	if !ok {
		return m, m.list.NewStatusMessage("Failed to mark task as 'Done'")
	}

	tsk := unmarshalTaskListItem(t)
	tsk.Status = int64(task.Done)

	_, err := m.TaskService.Update(tsk)
	if err != nil {
		return m, m.list.NewStatusMessage("Failed to mark task as 'Done'")
	}
	return m.updateListItems("Task marked as 'Done'")
}

func (m model) startTask(item list.Item) (tea.Model, tea.Cmd) {
	t, ok := item.(taskListItem)
	if !ok {
		return m, m.list.NewStatusMessage("Failed to mark task as 'Doing'")
	}

	tsk := unmarshalTaskListItem(t)
	tsk.Status = int64(task.Doing)

	_, err := m.TaskService.Update(tsk)
	if err != nil {
		return m, m.list.NewStatusMessage("Failed to mark task as 'Doing'")
	}
	return m.updateListItems("Task marked as 'Doing'")
}

func (m model) resetTask(item list.Item) (tea.Model, tea.Cmd) {
	t, ok := item.(taskListItem)
	if !ok {
		return m, m.list.NewStatusMessage("Failed to mark task as 'Todo'")
	}

	tsk := unmarshalTaskListItem(t)
	tsk.Status = int64(task.ToDo)

	_, err := m.TaskService.Update(tsk)
	if err != nil {
		return m, m.list.NewStatusMessage("Failed to mark task as 'Todo'")
	}
	return m.updateListItems("Task marked as 'Todo'")
}

func (m model) toggleTaskFavStatus(item list.Item) (tea.Model, tea.Cmd) {
	var verb string
	t, ok := item.(taskListItem)
	if !ok {
		return m, m.list.NewStatusMessage("Failed to mark task as 'Favorite'")
	}

	var favStatus task.FavoriteFlag
	switch t.IsFav {
	case int(task.IsFavorite):
		favStatus = task.IsNotFavorite
		verb = "unmarked"
	case int(task.IsNotFavorite):
		favStatus = task.IsFavorite
		verb = "marked"
	}

	tsk := unmarshalTaskListItem(t)
	tsk.IsFavorite = int(favStatus)

	_, err := m.TaskService.Update(tsk)
	if err != nil {
		return m, m.list.NewStatusMessage("Failed to mark task as 'Favorite'")
	}
	return m.updateListItems(fmt.Sprintf("Task %s as 'Favorite'", verb))
}

func (m model) listTitle() string {
	return category.GetCategoryNameFromID(m.Categories, m.SelectedCategoryID)
}
