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

	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/modules/category"
	"github.com/cgoesche/willdo/internal/modules/task"
	"github.com/spf13/cobra"
)

var (
	deleteAll bool

	deleteCmd = &cobra.Command{
		Use:   "delete [ category | task ]",
		Short: "Delete a category or task",
		Long: `There is not much more to say about this or 
are you looking for the entire commit history ?`,
		RunE: func(cmd *cobra.Command, args []string) error {
			db := database.New(conf.Database)
			catService := category.NewService(db)
			taskService := task.NewService(db)

			if len(args) <= 0 {
				return fmt.Errorf("missing argument 'category' or 'task'")
			}

			switch args[0] {
			case "task":
				if err := handleTasks(*taskService, *catService); err != nil {
					return err
				}
			case "category":
				if err := handleCategories(*taskService, *catService); err != nil {
					return err
				}
			default:
				return fmt.Errorf("missing argument 'category' or 'task'")
			}
			return nil
		},
	}
)

func init() {
	deleteCmd.Flags().StringVarP(&categoryName, "category", "c", "", "category name")
	deleteCmd.Flags().Int64VarP(&taskID, "task", "t", 0, "task ID")
	deleteCmd.Flags().BoolVarP(&deleteAll, "all", "a", false, "delete all categories/tasks")
}

func handleTasks(ts task.Service, cs category.Service) error {
	var err error
	var id int64
	if deleteAll {
		if err = ts.DeleteAll(); err != nil {
			return err
		}
		fmt.Printf("All tasks deleted\n")
	} else if categoryName != "" {
		cats, err := cs.GetAll()
		if err != nil {
			return fmt.Errorf("failed to find any categories in the database, %v", err)
		}

		categoryID = category.GetCategoryIDFromName(cats, categoryName)
		if categoryID <= 0 {
			return fmt.Errorf("failed to get category ID for name %s", categoryName)
		}

		if _, err = ts.DeleteAllByCategory(categoryID); err != nil {
			return err
		}
		fmt.Printf("Tasks of category '%s' deleted\n", categoryName)
	} else {
		if id, err = ts.Delete(taskID); err != nil {
			return err
		}
		fmt.Printf("Task %d deleted\n", id)
	}
	return nil
}

func handleCategories(ts task.Service, cs category.Service) error {
	var err error
	if deleteAll {
		if err = cs.DeleteAll(); err != nil {
			return err
		}

		if err = ts.DeleteAll(); err != nil {
			return err
		}
		fmt.Printf("All categories deleted\n")
	} else {
		cats, err := cs.GetAll()
		if err != nil {
			return fmt.Errorf("failed to find any categories in the database, %v", err)
		}

		categoryID = category.GetCategoryIDFromName(cats, categoryName)
		if categoryID <= 0 {
			return fmt.Errorf("failed to get category ID for name %s", categoryName)
		}

		if _, err = cs.Delete(categoryID); err != nil {
			return err
		}

		if _, err = ts.DeleteAllByCategory(categoryID); err != nil {
			return err
		}
		fmt.Printf("Category '%s' deleted\n", categoryName)
	}
	return nil
}
