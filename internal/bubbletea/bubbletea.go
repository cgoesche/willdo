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
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(m model) {
	var l []list.Item
	var err error
	var d *taskItemDelegate
	var defaultList list.Model

	d = newTaskItemDelegate()
	d.categories = m.Categories
	defaultList = list.New([]list.Item{}, d, 0, 0)
	defaultList.SetStatusBarItemName("task", "tasks")
	m.list = defaultList

	if m.ShowAllTasks {
		l, err = m.getAllTaskListItems()
		m.list.Title = "All tasks"
		d.showCategory = true
	} else {
		l, err = m.getTaskListItemsByCategory(m.SelectedCategoryID)
		m.list.Title = m.listTitle()
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	m.cachedItems = l

	m.list.SetItems(m.cachedItems)
	m.list.Styles = styles.DefaultStyles()
	m.list.StatusMessageLifetime = 3 * time.Second
	m.list.SetShowFilter(false)
	m.list.SetFilteringEnabled(false)
	m.list.InfiniteScrolling = true
	m.list.AdditionalShortHelpKeys = m.KeyMap.ShortHelpKeys
	m.list.AdditionalFullHelpKeys = m.KeyMap.FullHelpKeys

	m.details.selectedItem = m.list.SelectedItem()

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
