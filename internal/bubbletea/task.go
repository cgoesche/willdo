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
	"io"
	"strings"

	"github.com/cgoesche/willdo/internal/bubbletea/styles"
	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/models"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type taskListItem struct {
	ID   int64
	Tit  string
	Desc string
	Stat int64
	Prio int64
	Cat  int64
}

func (i taskListItem) Identity() int64     { return i.ID }
func (i taskListItem) Title() string       { return i.Tit }
func (i taskListItem) Description() string { return i.Desc }
func (i taskListItem) Status() int64       { return i.Stat }
func (i taskListItem) Priority() int64     { return i.Prio }
func (i taskListItem) Category() int64     { return i.Cat }
func (i taskListItem) FilterValue() string { return i.Tit }

type taskItemDelegate struct {
	height  int
	spacing int
}

func (d taskItemDelegate) Height() int {
	return d.height
}

func (d taskItemDelegate) Spacing() int {
	return d.spacing
}

func (d taskItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d taskItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(taskListItem)
	if !ok {
		return
	}

	id := styles.TaskIdentityStyle.Render(fmt.Sprintf("%d", i.Identity()))
	str := fmt.Sprintf("%s.  %s  %s", id, styles.RenderStatusIcon(models.Status(i.Status())), i.Title())

	fn := styles.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.SelectedItemStyle.Render(">" + strings.Join(s, ""))
		}
	}

	fmt.Fprint(w, fn(str))
}

func (d *taskItemDelegate) SetHeight(h int) {
	d.height = h
}

func (d *taskItemDelegate) SetSapcing(s int) {
	d.spacing = s
}

func newTaskItemDelegate() *taskItemDelegate {
	return &taskItemDelegate{
		spacing: 0,
		height:  1,
	}
}

func getTaskListItems(c *database.Client) ([]list.Item, error) {
	var l []list.Item

	tasks, err := c.QueryAllTasks()
	if err != nil {
		return nil, err
	}

	for _, v := range tasks {
		var i taskListItem
		i.ID = v.ID
		i.Tit = v.Title
		i.Stat = v.Status
		i.Desc = v.Description

		l = append(l, i)
	}
	return l, nil
}

func getTaskListItemsByCategory(c *database.Client, id int64) ([]list.Item, error) {
	var l []list.Item

	tasks, err := c.QueryTasksFromCategory(id)
	if err != nil {
		return nil, err
	}

	for _, v := range tasks {
		var i taskListItem
		i.ID = v.ID
		i.Tit = v.Title
		i.Stat = v.Status
		i.Desc = v.Description

		l = append(l, i)
	}
	return l, nil
}

func (m model) getTaskItemID() (int64, error) {
	i := m.lists[m.listIndex].Index()
	items := m.lists[m.listIndex].Items()

	task, ok := items[i].(taskListItem)
	if !ok {
		return -1, fmt.Errorf("failed to get task item ID")
	}
	return task.Identity(), nil
}

func (m model) deleteTask() (tea.Model, tea.Cmd) {
	id, err := m.getTaskItemID()
	if err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage(fmt.Sprintf("%v", err))
	}

	err = m.dbClient.DeleteTask(id)
	if err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage("Failed to delete task")
	}

	if err = m.updateListItems(); err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage("Task deleted but failed to update task.")
	}

	return m, m.lists[m.listIndex].NewStatusMessage("Deleted task")
}

func (m model) completeTask() (tea.Model, tea.Cmd) {
	id, err := m.getTaskItemID()
	if err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage(fmt.Sprintf("%v", err))
	}

	_, err = m.dbClient.CompleteTask(int(id))
	if err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage("Failed to mark task as completed")
	}

	if err = m.updateListItems(); err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage("Task marked as completed but failed to update task list.")
	}

	return m, m.lists[m.listIndex].NewStatusMessage("Task marked as 'Done'")
}

func (m model) startTask() (tea.Model, tea.Cmd) {
	id, err := m.getTaskItemID()
	if err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage(fmt.Sprintf("%v", err))
	}

	_, err = m.dbClient.StartTask(int(id))
	if err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage("Failed to mark task")
	}

	if err = m.updateListItems(); err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage("Task started but failed to update task list.")
	}

	return m, m.lists[m.listIndex].NewStatusMessage("Task marked as 'Doing'")
}

func (m model) resetTask() (tea.Model, tea.Cmd) {
	id, err := m.getTaskItemID()
	if err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage(fmt.Sprintf("%v", err))
	}

	_, err = m.dbClient.ResetTask(int(id))
	if err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage("Failed to mark task as todo")
	}

	if err = m.updateListItems(); err != nil {
		return m, m.lists[m.listIndex].NewStatusMessage("Task reset but failed to update task list.")
	}

	return m, m.lists[m.listIndex].NewStatusMessage("Task marked as 'Todo'")
}
