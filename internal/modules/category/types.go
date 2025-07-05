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
package category

type Category struct {
	ID          int64
	Name        string
	Description string
}

type Categories []Category
type CategoryNameToIDMap map[string]int64
type CategoryIDToNameMap map[int64]string

func GetCategoryNameFromID(cats Categories, id int64) string {
	for _, v := range cats {
		if id == v.ID {
			return v.Name
		}
	}
	return "N/A"
}

func GetCategoryIDFromName(cats Categories, n string) int64 {
	for _, v := range cats {
		if n == v.Name {
			return v.ID
		}
	}
	return 0
}

func NewCategoryNameToIDMap(cats Categories) CategoryNameToIDMap {
	m := make(CategoryNameToIDMap, len(cats))
	for _, v := range cats {
		m[v.Name] = v.ID
	}
	return m
}

func NewCategoryIDToNameMap(cats Categories) CategoryIDToNameMap {
	m := make(CategoryIDToNameMap, len(cats))
	for _, v := range cats {
		m[v.ID] = v.Name
	}
	return m
}
