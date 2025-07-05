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
	"github.com/cgoesche/willdo/internal/modules/category"
	"github.com/cgoesche/willdo/internal/modules/task"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

const (
	ellipsis = "..."
	bullet
)

type taskListItem struct {
	ID    int64
	Tit   string
	Desc  string
	Stat  int64
	Prio  int64
	Cat   int64
	IsFav int
}

func (i taskListItem) Identity() int64     { return i.ID }
func (i taskListItem) Title() string       { return i.Tit }
func (i taskListItem) Description() string { return i.Desc }
func (i taskListItem) Status() int64       { return i.Stat }
func (i taskListItem) Priority() int64     { return i.Prio }
func (i taskListItem) Category() int64     { return i.Cat }
func (i taskListItem) IsFavorite() int     { return i.IsFav }
func (i taskListItem) FilterValue() string { return i.Tit }

type taskItemDelegate struct {
	height       int
	spacing      int
	categories   category.Categories
	showCategory bool
}

func newTaskItemDelegate() *taskItemDelegate {
	return &taskItemDelegate{
		spacing:      0,
		height:       1,
		categories:   nil,
		showCategory: false,
	}
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
	t, ok := listItem.(taskListItem)
	if !ok {
		return
	}

	var str string
	num := fmt.Sprintf("%   3d.", index+1)
	num = styles.SubtleStyle.Render(num)
	title := ansi.Truncate(t.Title(), m.Width()/2, ellipsis)
	statusIcon := styles.RenderStatusIcon(task.Status(t.Stat))

	if t.Priority() == int64(task.High) && t.Status() != int64(task.Done) {
		title = styles.HighPriorityStyle.Render(title)
	}

	if t.Status() == int64(task.Done) {
		title = styles.SubtleStyle.Render(title)
	}

	str = fmt.Sprintf("%s %s %s", num, statusIcon, title)

	if t.IsFavorite() == task.IsFavorite {
		str += styles.FavoriteIconStyle.Render(" " + task.FavoriteIcon)
	}

	if t.Description() != "" {
		str += styles.NoteIndicatorStyle.Render(" " + task.NoteIndicatorIcon)
	}

	if d.showCategory {
		str += styles.TaskCategoryNameStyle.Render("[" + category.GetCategoryNameFromID(d.categories, t.Cat) + "]")
	}

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

func (d *taskItemDelegate) SetSpacing(s int) {
	d.spacing = s
}
