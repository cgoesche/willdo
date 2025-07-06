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
package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	ClearFilter     key.Binding
	CompleteTask    key.Binding
	DeleteTask      key.Binding
	FilterToDo      key.Binding
	FilterDoing     key.Binding
	FilterDone      key.Binding
	FilterFav       key.Binding
	NextCategory    key.Binding
	PrevCategory    key.Binding
	ResetTask       key.Binding
	RefreshList     key.Binding
	StartTask       key.Binding
	ToggleFavStatus key.Binding
	ToggleDetails   key.Binding
	ToggleStats     key.Binding
}

var DefaultKeyMap = KeyMap{
	ClearFilter: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("^r", "reset filter"),
	),
	CompleteTask: key.NewBinding(
		key.WithKeys("c", " "),
		key.WithHelp("c/␣", "complete task"),
	),
	DeleteTask: key.NewBinding(
		key.WithKeys("d", "backspace"),
		key.WithHelp("d", "delete task"),
	),
	FilterToDo: key.NewBinding(
		key.WithKeys("ctrl+t"),
		key.WithHelp("^t", "filter todo tasks"),
	),
	FilterDoing: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("^s", "filter doing tasks"),
	),
	FilterDone: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("^d", "filter done tasks"),
	),
	FilterFav: key.NewBinding(
		key.WithKeys("ctrl+f"),
		key.WithHelp("^f", "filter favorites"),
	),
	NextCategory: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next category"),
	),
	PrevCategory: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "previous category"),
	),
	ResetTask: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "mark as todo"),
	),
	RefreshList: key.NewBinding(
		key.WithKeys("alt+r"),
		key.WithHelp("alt-r", "refresh list"),
	),
	StartTask: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start task"),
	),
	ToggleFavStatus: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "(un)mark as favorite"),
	),
	ToggleDetails: key.NewBinding(
		key.WithKeys("alt+d"),
		key.WithHelp("alt-d", "show/hide details"),
	),
	ToggleStats: key.NewBinding(
		key.WithKeys("alt+s"),
		key.WithHelp("alt-s", "show/hide stats"),
	),
}

func (k KeyMap) ShortHelpKeys() []key.Binding {
	kb := []key.Binding{
		k.CompleteTask,
		k.DeleteTask,
		k.ResetTask,
		k.StartTask,
	}
	return kb
}

func (k KeyMap) FullHelpKeys() []key.Binding {
	kb := []key.Binding{
		k.CompleteTask,
		k.DeleteTask,
		k.FilterToDo,
		k.FilterDoing,
		k.FilterDone,
		k.FilterFav,
		k.NextCategory,
		k.PrevCategory,
		k.ResetTask,
		k.RefreshList,
		k.StartTask,
		k.ToggleFavStatus,
		k.ToggleDetails,
		k.ToggleStats,
	}
	return kb
}
