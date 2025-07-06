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
	lists        []list.Model
	selectedList int
	stats        statsModel
	details      detailsModel

	CategoryService  *category.Service
	TaskService      *task.Service
	Categories       category.Categories
	CatNameToIDMap   category.CategoryNameToIDMap
	CatIDToNameMap   category.CategoryIDToNameMap
	SelectedCategory int64

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
	item := m.lists[m.selectedList].SelectedItem()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.ToggleDetails):
			m.ShowDetails = !m.ShowDetails
			return m, m.lists[m.selectedList].NewStatusMessage("Toggled 'Details' section")
		case key.Matches(msg, m.KeyMap.ToggleStats):
			m.ShowStats = !m.ShowStats
			return m, m.lists[m.selectedList].NewStatusMessage("Toggled 'Statistics' section")
		case key.Matches(msg, m.KeyMap.RefreshList):
			return m.updateListItems()
		case key.Matches(msg, m.KeyMap.DeleteTask):
			return m.deleteTask(item)
		case key.Matches(msg, m.KeyMap.CompleteTask):
			return m.completeTask(item)
		case key.Matches(msg, m.KeyMap.ToggleFavStatus):
			return m.toggleTaskFavStatus(item)
		case key.Matches(msg, m.KeyMap.StartTask):
			return m.startTask(item)
		case key.Matches(msg, m.KeyMap.ResetTask):
			return m.resetTask(item)
		case key.Matches(msg, m.KeyMap.FilterFav):
			m.IsFiltering = true
			m.FilterValue = task.IsFavorite
			return m.updateListItems()
		case key.Matches(msg, m.KeyMap.ClearFilter):
			m.IsFiltering = false
			m.FilterValue = nil
			return m.updateListItems()
		case key.Matches(msg, m.KeyMap.FilterToDo):
			m.IsFiltering = true
			m.FilterValue = task.ToDo
			return m.updateListItems()
		case key.Matches(msg, m.KeyMap.FilterDoing):
			m.IsFiltering = true
			m.FilterValue = task.Doing
			return m.updateListItems()
		case key.Matches(msg, m.KeyMap.FilterDone):
			m.IsFiltering = true
			m.FilterValue = task.Done
			return m.updateListItems()
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

func (m model) getTaskListItemsWithFilter() ([]list.Item, error) {
	var l []list.Item
	var catID = m.SelectedCategory
	if m.ShowAllTasks {
		catID = -1
	}
	tasks, err := m.TaskService.GetAllWithFilter(catID, m.FilterValue)
	if err != nil {
		return nil, err
	}

	l = marshalTaskListItems(tasks)
	return l, nil
}

func (m *model) updateListItems() (tea.Model, tea.Cmd) {
	var l []list.Item
	var err error

	switch {
	case m.ShowAllTasks && !m.IsFiltering:
		l, err = m.getAllTaskListItems()
	case m.IsFiltering:
		l, err = m.getTaskListItemsWithFilter()
	default:
		l, err = m.getTaskListItemsByCategory(m.SelectedCategory)
	}

	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage(fmt.Sprintf("Failed to refresh list, %v", err))
	}
	m.lists[m.selectedList].SetItems(l)
	return m, m.lists[m.selectedList].NewStatusMessage("List refreshed")
}

func (m model) deleteTask(item list.Item) (tea.Model, tea.Cmd) {
	task, ok := item.(taskListItem)
	if !ok {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to delete task")
	}

	_, err := m.TaskService.Delete(task.ID)
	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to delete task")
	}
	return m.updateListItems()
}

func (m model) completeTask(item list.Item) (tea.Model, tea.Cmd) {
	t, ok := item.(taskListItem)
	if !ok {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Done'")
	}

	tsk := unmarshalTaskListItem(t)
	tsk.Status = int64(task.Done)

	_, err := m.TaskService.Update(tsk)
	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Done'")
	}
	return m.updateListItems()
}

func (m model) startTask(item list.Item) (tea.Model, tea.Cmd) {
	t, ok := item.(taskListItem)
	if !ok {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Doing'")
	}

	tsk := unmarshalTaskListItem(t)
	tsk.Status = int64(task.Doing)

	_, err := m.TaskService.Update(tsk)
	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Doing'")
	}
	return m.updateListItems()
}

func (m model) resetTask(item list.Item) (tea.Model, tea.Cmd) {
	t, ok := item.(taskListItem)
	if !ok {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Todo'")
	}

	tsk := unmarshalTaskListItem(t)
	tsk.Status = int64(task.ToDo)

	_, err := m.TaskService.Update(tsk)
	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Todo'")
	}
	return m.updateListItems()
}

func (m model) toggleTaskFavStatus(item list.Item) (tea.Model, tea.Cmd) {
	t, ok := item.(taskListItem)
	if !ok {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Favorite'")
	}

	var favStatus task.FavoriteFlag
	switch t.IsFav {
	case int(task.IsFavorite):
		favStatus = task.IsNotFavorite
	case int(task.IsNotFavorite):
		favStatus = task.IsFavorite
	}

	tsk := unmarshalTaskListItem(t)
	tsk.IsFavorite = int(favStatus)

	_, err := m.TaskService.Update(tsk)
	if err != nil {
		return m, m.lists[m.selectedList].NewStatusMessage("Failed to mark task as 'Favorite'")
	}
	return m.updateListItems()
}

func (m model) listTitle() string {
	return category.GetCategoryNameFromID(m.Categories, m.SelectedCategory)
}
