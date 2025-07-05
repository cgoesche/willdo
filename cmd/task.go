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
	"github.com/cgoesche/willdo/internal/modules/task"
	"github.com/spf13/cobra"
)

var (
	t = &task.Task{}

	taskCmd = &cobra.Command{
		Use:   "task",
		Short: "Add a task",
		Long: `There is not much more to say about this or 
are you looking for the entire commit history ?`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db := database.New(conf.Database)

			catService := category.NewService(db)
			taskService := task.NewService(db)

			err := catService.InitRepo()
			if err != nil {
				return fmt.Errorf("failed to init category service repository, %v", err)
			}
			err = taskService.InitRepo()
			if err != nil {
				return fmt.Errorf("failed to init task service repository, %v", err)
			}

			if strings.TrimSpace(categoryName) == "" {
				return fmt.Errorf("category name cannot be an empty string")
			}

			cats, err := catService.GetAll()
			if err != nil {
				return fmt.Errorf("failed to find any categories in the database, %v", err)
			}
			var categoryID int64
			if categoryName != "" {
				categoryID = category.GetCategoryIDFromName(cats, categoryName)
				if categoryID == 0 {
					return fmt.Errorf("no category found with name '%s'", categoryName)
				}
			}
			t.Category = categoryID

			var id int64
			if id, err = taskService.Create(*t); err != nil {
				return fmt.Errorf("failed to add task: %v", err)
			}

			fmt.Printf("Task %d added!\n", id)
			return nil
		},
	}
)

func init() {
	taskCmd.Flags().StringVarP(&categoryName, "category", "c", "", "task category")
	taskCmd.Flags().StringVarP(&t.Title, "title", "t", "", "task title")
	taskCmd.Flags().StringVarP(&t.Description, "description", "d", "", "task description")
	taskCmd.Flags().Int64VarP(&t.Status, "status", "s", 0, "task status (e.g. 0, 1, or 2)")
	taskCmd.Flags().Int64VarP(&t.Priority, "priority", "p", 0, "task priority (e.g. 0, 1, or 2)")
	taskCmd.Flags().IntVarP(&t.IsFavorite, "favorite", "f", 0, "Mark task as favorite")

	taskCmd.MarkFlagsOneRequired("title", "category")
	taskCmd.MarkFlagsRequiredTogether("title", "category")
}
