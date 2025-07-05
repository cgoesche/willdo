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
	"github.com/cgoesche/willdo/internal/modules/task"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	NormalColor       = lipgloss.Color("#EEEEEE")
	LightSubtleColor  = lipgloss.Color("#8c8c8c")
	SubtleColor       = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#4d4d4d"}
	SpecialColor      = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	HighlightColor    = lipgloss.Color("#EF56f4")
	NoticeColor       = lipgloss.Color("111")
	WarningColor      = lipgloss.Color("#FF4E4E")
	HighPriorityStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF344")).Underline(true)

	DocStyle                    = lipgloss.NewStyle().Margin(1, 2)
	TitleBarStyle               = lipgloss.NewStyle().PaddingTop(1)
	ListTitleStyle              = lipgloss.NewStyle().Foreground(SpecialColor).Underline(true).Bold(true)
	TaskDetailSectionTitleStyle = lipgloss.NewStyle().Foreground(SubtleColor).PaddingBottom(1).Underline(true).Bold(true)
	TaskCategoryNameStyle       = lipgloss.NewStyle().Foreground(SubtleColor).PaddingLeft(1)
	ItemStyle                   = lipgloss.NewStyle().PaddingLeft(1)
	SelectedItemStyle           = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("170")).Bold(true)
	TaskIdentityStyle           = lipgloss.NewStyle().PaddingLeft(2).Foreground(SubtleColor)
	StatusBarStyle              = lipgloss.NewStyle().Foreground(SubtleColor).MarginBottom(1)
	PaginationStyle             = list.DefaultStyles().PaginationStyle.PaddingLeft(1)
	HelpStyle                   = lipgloss.NewStyle().PaddingLeft(0).Foreground(SubtleColor).MarginBottom(1)
	QuitTextStyle               = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	FavoriteIconStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF344"))
	NoteIndicatorStyle          = lipgloss.NewStyle().Foreground(NoticeColor)
	TodoStyle                   = lipgloss.NewStyle().Foreground(HighlightColor)
	DoingStyle                  = lipgloss.NewStyle().Foreground(NoticeColor)
	DoneStyle                   = lipgloss.NewStyle().Foreground(SpecialColor)
	SubtleStyle                 = lipgloss.NewStyle().Foreground(SubtleColor)
)

func DefaultStyles() list.Styles {
	s := list.DefaultStyles()
	s.TitleBar = TitleBarStyle
	s.Title = ListTitleStyle
	s.StatusBar = StatusBarStyle
	s.HelpStyle = HelpStyle
	s.PaginationStyle = PaginationStyle
	return s
}

func RenderStatusIcon(s task.Status) string {
	var style lipgloss.Style

	switch s {
	case task.ToDo:
		style = TodoStyle
	case task.Doing:
		style = DoingStyle
	case task.Done:
		style = DoneStyle

	}
	return style.Render(task.StatusMap[task.Status(s)])
}
