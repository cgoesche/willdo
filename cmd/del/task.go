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

	"github.com/cgoesche/willdo/internal/config"
	"github.com/cgoesche/willdo/internal/database"
	"github.com/spf13/cobra"
)

var (
	taskID         int64
	taskCategory   int64
	deleteAllTasks bool

	taskCmd = &cobra.Command{
		Use:   "task",
		Short: "Delete a task",
		Long: `There is not much more to say about this or 
are you looking for the entire commit history ?`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client := database.NewClient()
			err := client.InitDB(config.SetDefault().Database.Path)
			if err != nil {
				return err
			}

			if deleteAllTasks {
				err = client.DeleteAllTasks()
				if err != nil {
					return err
				}
				fmt.Printf("All tasks deleted successfully!\n")
				return nil
			}

			if err = client.DeleteTask(taskID); err != nil {
				return err
			}
			fmt.Printf("Task (ID: %d) deleted successfully!\n", taskID)
			return nil
		},
	}
)

func init() {
	taskCmd.Flags().Int64VarP(&taskID, "id", "i", 0, "Task ID")
	taskCmd.Flags().Int64VarP(&taskCategory, "category", "c", 0, "Task category")
	taskCmd.Flags().BoolVarP(&deleteAllTasks, "all", "a", false, "Delete all tasks")

	taskCmd.MarkFlagsOneRequired("id", "all")
	taskCmd.MarkFlagsMutuallyExclusive("id", "all")
}
