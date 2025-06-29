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
package add

import (
	"fmt"

	"github.com/cgoesche/willdo/internal/config"
	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/models"
	"github.com/spf13/cobra"
)

var (
	cat = &models.Category{}

	categoryCmd = &cobra.Command{
		Use:   "category",
		Short: "Add a category",
		Long: `There is not much more to say about this or 
are you looking for the entire commit history ?`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := database.NewClient()
			err := client.InitDB(config.SetDefault().Database.Path)
			if err != nil {
				return err
			}

			if err = addCategory(client); err != nil {
				return fmt.Errorf("failed to add category: %v", err)
			}

			fmt.Printf("Category '%s' added successfully!\n", cat.Name)
			return nil
		},
	}
)

func init() {
	categoryCmd.Flags().StringVarP(&cat.Name, "name", "n", "", "Category name")
	categoryCmd.Flags().StringVarP(&cat.Description, "description", "d", "", "Category description")

	categoryCmd.MarkFlagRequired("name")
}

func addCategory(client *database.Client) error {
	cat := &models.Category{
		Name:        cat.Name,
		Description: cat.Description,
	}
	_, err := client.InsertRow(cat)
	if err != nil {
		return err
	}

	return nil
}
