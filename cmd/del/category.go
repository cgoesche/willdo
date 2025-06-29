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
package del

import (
	"fmt"

	"github.com/cgoesche/willdo/internal/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	categoryID   int64
	deleteAllCat bool

	categoryCmd = &cobra.Command{
		Use:   "category",
		Short: "Delete a category",
		Long: `There is not much more to say about this or 
are you looking for the entire commit history ?`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := viper.GetViper().GetString("database.path")
			client := database.NewClient()
			err := client.InitDB(path)
			if err != nil {
				return err
			}

			if deleteAllCat {
				err := client.DeleteAllCategories()
				if err != nil {
					return fmt.Errorf("failed to delete all categories, %v", err)
				}
				fmt.Printf("All categories deleted successfully!\n")
				return nil
			}

			if err = client.DeleteCategory(categoryID); err != nil {
				return err
			}

			fmt.Printf("Category (ID: %d) deleted successfully!\n", categoryID)
			return nil
		},
	}
)

func init() {
	categoryCmd.Flags().Int64VarP(&categoryID, "id", "i", 0, "Category ID")
	categoryCmd.Flags().BoolVarP(&deleteAllCat, "all", "a", false, "Delete all categories")

	categoryCmd.MarkFlagsOneRequired("id", "all")
	categoryCmd.MarkFlagsMutuallyExclusive("id", "all")
}
