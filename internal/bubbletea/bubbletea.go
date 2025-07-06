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

	"github.com/cgoesche/willdo/internal/bubbletea/keys"
	"github.com/cgoesche/willdo/internal/bubbletea/styles"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(m model) {
	m.KeyMap = keys.DefaultKeyMap
	d := newTaskItemDelegate()
	d.categories = m.Categories
	defaultList := list.New([]list.Item{}, d, 0, 0)
	defaultList.SetStatusBarItemName("task", "tasks")
	m.lists = []list.Model{defaultList, defaultList}

	var l []list.Item
	var err error
	if m.ShowAllTasks {
		l, err = m.getAllTaskListItems()
		m.lists[m.selectedList].Title = "All tasks"
		d.showCategory = true
	} else {
		l, err = m.getTaskListItemsByCategory(m.SelectedCategory)
		m.lists[m.selectedList].Title = m.listTitle()
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	m.lists[m.selectedList].SetItems(l)
	m.lists[m.selectedList].Styles = styles.DefaultStyles()
	m.lists[m.selectedList].StatusMessageLifetime = 3 * time.Second
	m.lists[m.selectedList].SetShowFilter(false)
	m.lists[m.selectedList].SetFilteringEnabled(false)
	m.lists[m.selectedList].InfiniteScrolling = true
	m.lists[m.selectedList].AdditionalShortHelpKeys = m.KeyMap.ShortHelpKeys
	m.lists[m.selectedList].AdditionalFullHelpKeys = m.KeyMap.FullHelpKeys

	m.details.selectedItem = m.lists[m.selectedList].SelectedItem()

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
