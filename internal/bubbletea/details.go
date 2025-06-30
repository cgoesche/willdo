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

	"github.com/cgoesche/willdo/internal/bubbletea/styles"
)

func (m model) RenderTaskDetailsSection() string {
	i := m.lists[m.listIndex].Index()
	items := m.lists[m.listIndex].Items()

	item, ok := items[i].(taskListItem)
	if !ok {
		return ""
	}

	sectionTitle := styles.TaskDetailSectionTitleStyle.Render("Details")

	section := fmt.Sprintf("%s\nID: %d\nTitle: %s\nDescription: %s\n\n", sectionTitle, item.Identity(), item.Title(), item.Description())
	return section
}
