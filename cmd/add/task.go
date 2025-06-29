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

	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	t = &models.Task{}

	taskCmd = &cobra.Command{
		Use:   "task",
		Short: "Add a task",
		Long: `There is not much more to say about this or 
are you looking for the entire commit history ?`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := viper.GetViper().GetString("database.path")
			client := database.NewClient()

			err := client.InitDB(path)
			if err != nil {
				return err
			}

			if err = addTask(client); err != nil {
				return fmt.Errorf("failed to add task: %v", err)
			}

			fmt.Printf("Task '%s' added successfully!\n", t.Title)
			return nil
		},
	}
)

func init() {
	taskCmd.Flags().StringVarP(&t.Title, "title", "t", "", "Task title")
	taskCmd.Flags().StringVarP(&t.Description, "description", "d", "", "Task description")
	taskCmd.Flags().Int64VarP(&t.Status, "status", "s", 0, "Task status")
	taskCmd.Flags().Int64VarP(&t.Priority, "priority", "p", 0, "Task priority")
	taskCmd.Flags().Int64VarP(&t.Category, "category", "c", 1, "Task category")

	taskCmd.MarkFlagRequired("title")
}

func addTask(c *database.Client) error {
	if t.Status > int64(models.Done) || t.Status < int64(models.ToDo) {
		return fmt.Errorf("invalid status value")
	}

	if t.Priority < int64(models.Low) || t.Status > int64(models.Urgent) {
		return fmt.Errorf("invalid priority value")
	}

	task := &models.Task{
		Title:       t.Title,
		Description: t.Description,
		Priority:    t.Priority,
		Status:      t.Status,
		Category:    t.Category,
	}
	_, err := c.InsertRow(task)
	if err != nil {
		return err
	}
	return nil
}
