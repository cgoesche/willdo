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
	"fmt"
	"strings"

	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/modules/category"
	"github.com/cgoesche/willdo/internal/modules/task"
	"github.com/spf13/cobra"
)

var (
	taskID          int64
	taskTitle       string
	taskDescription string
	taskStatus      int64
	taskPriority    int64
	taskIsFavorite  int
	categoryID      int64
	categoryName    string

	taskCmd = &cobra.Command{
		Use:   "task",
		Short: "Edit a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			db := database.New(conf.Database)
			taskService := task.NewService(db)
			catService := category.NewService(db)

			cats, err := catService.GetAll()
			if err != nil {
				return fmt.Errorf("failed to find any categories in the database, %v", err)
			}

			tsk, err := taskService.GetById(taskID)
			if err != nil {
				return fmt.Errorf("could not find a task with ID %d", taskID)
			}

			if len(strings.TrimSpace(taskTitle)) != 0 {
				tsk.Title = taskTitle
			}
			if taskDescription != "" {
				tsk.Description = taskDescription
			}
			if taskStatus != tsk.Status && taskStatus >= int64(task.ToDo) && taskStatus <= int64(task.Done) {
				tsk.Status = taskStatus
			}
			if taskPriority != tsk.Priority && taskPriority >= int64(task.Low) {
				tsk.Priority = taskPriority
			}
			if int64(taskIsFavorite) != int64(tsk.IsFavorite) && int64(taskIsFavorite) == int64(task.IsFavorite) || int64(taskIsFavorite) <= int64(task.IsNotFavorite) {
				tsk.IsFavorite = taskIsFavorite
			}
			if len(strings.TrimSpace(categoryName)) != 0 {
				categoryID = category.GetCategoryIDFromName(cats, categoryName)
				if categoryID == 0 {
					return fmt.Errorf("no category found with name '%s'", categoryName)
				}
				tsk.Category = categoryID
			}

			id, err := taskService.Update(tsk)
			if err != nil {
				return fmt.Errorf("failed to update task %d, %v", id, err)
			}
			fmt.Printf("Task %d updated!\n", taskID)
			return nil
		},
	}
)

func init() {
	taskCmd.Flags().Int64VarP(&taskID, "id", "i", 0, "task ID")
	taskCmd.Flags().StringVarP(&categoryName, "category", "c", "", "new task category")
	taskCmd.Flags().StringVarP(&taskTitle, "title", "t", "", "new task title")
	taskCmd.Flags().StringVarP(&taskDescription, "description", "d", "", "new task description")
	taskCmd.Flags().Int64VarP(&taskStatus, "status", "s", -1, "new task status (e.g. 0, 1, or 2)")
	taskCmd.Flags().Int64VarP(&taskPriority, "priority", "p", -1, "new task priority (e.g. 0, 1, or 2)")
	taskCmd.Flags().IntVarP(&taskIsFavorite, "favorite", "f", -1, "mark task as favorite")

	taskCmd.MarkFlagRequired("id")
}
