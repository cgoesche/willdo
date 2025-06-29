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
package styles

import (
	"github.com/cgoesche/willdo/internal/models"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(1)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
	TaskIdentityStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(subtle)
	PaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(1)
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(2).PaddingBottom(1)
	QuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	normal    = lipgloss.Color("#EEEEEE")
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	Special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	Highlight = lipgloss.Color("#EF56f4")
	Notice    = lipgloss.Color("111")
)

func DefaultStyles() list.Styles {
	s := list.DefaultStyles()

	s.TitleBar = lipgloss.NewStyle().Padding(0, 0, 0, 0)

	s.Title = lipgloss.NewStyle().
		Foreground(Special).Underline(true).Bold(true)

	s.StatusBar = lipgloss.NewStyle().
		Foreground(normal)

	s.PaginationStyle = PaginationStyle

	return s
}

func RenderStatusIcon(s models.Status) string {
	var style lipgloss.Style

	switch s {
	case models.ToDo:
		style = lipgloss.NewStyle().Foreground(Highlight)
	case models.Doing:
		style = lipgloss.NewStyle().Foreground(Notice)
	case models.Done:
		style = lipgloss.NewStyle().Foreground(Special)

	}
	return style.Render(models.StatusMap[models.Status(s)])
}
