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
package models

type Priority int64

const (
	Low Priority = iota
	Medium
	High
	Urgent
)

type Status int64

const (
	ToDo Status = iota
	Doing
	Done
)

var (
	StatusMap map[Status]string = map[Status]string{
		ToDo:  "☐",
		Doing: "…",
		Done:  "✔",
	}

	FavoriteIcon  = "★"
	IsFavorite    = 1
	IsNotFavorite = 0

	NoteIndicatorIcon = "🛈"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	Priority    int64
	Status      int64
	Category    int64
	IsFavorite  int
}

type Tasks []Task
