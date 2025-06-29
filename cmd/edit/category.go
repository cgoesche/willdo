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
package edit

import (
	"github.com/cgoesche/willdo/internal/models"
	"github.com/spf13/cobra"
)

var (
	c = &models.Category{}

	categoryCmd = &cobra.Command{
		Use:   "category",
		Short: "Edit a category",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func init() {
	categoryCmd.Flags().Int64VarP(&c.ID, "id", "i", 0, "Category ID")
	categoryCmd.Flags().StringVarP(&c.Name, "name", "n", "", "Category name")
	categoryCmd.Flags().StringVarP(&c.Description, "description", "d", "", "Category description")

	categoryCmd.MarkFlagRequired("id")
}
