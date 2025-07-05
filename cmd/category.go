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
package cmd

import (
	"fmt"
	"strings"

	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/modules/category"
	"github.com/spf13/cobra"
)

var (
	cat = &category.Category{}

	categoryCmd = &cobra.Command{
		Use:   "category",
		Short: "Add a category",
		Long: `There is not much more to say about this or 
are you looking for the entire commit history ?`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db := database.New(conf.Database)
			catService := category.NewService(db)

			err := catService.InitRepo()
			if err != nil {
				return fmt.Errorf("failed to init category service repository, %v", err)
			}

			if strings.TrimSpace(cat.Name) == "" {
				return fmt.Errorf("category name cannot be an empty string")
			}

			var id int64
			if id, err = catService.Create(*cat); err != nil {
				return fmt.Errorf("failed to add category: %v", err)
			}

			fmt.Printf("Category %d added!\n", id)
			return nil
		},
	}
)

func init() {
	categoryCmd.Flags().StringVarP(&cat.Name, "name", "n", "", "Category name")
	categoryCmd.Flags().StringVarP(&cat.Description, "description", "d", "", "Category description")
	categoryCmd.MarkFlagRequired("name")
}
