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
	"strconv"

	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/modules/task"
	"github.com/spf13/cobra"
)

var (
	resetCmd = &cobra.Command{
		Use:   "reset [ID]",
		Short: "Reset a task's status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db := database.New(conf.Database)
			taskService := task.NewService(db)

			err := taskService.InitRepo()
			if err != nil {
				return fmt.Errorf("failed to init task service repository, %v", err)
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task ID")
			}

			t, err := taskService.GetById(int64(id))
			if err != nil {
				return fmt.Errorf("failed to get task %d, %v", id, err)
			}

			t.Status = int64(task.ToDo)
			ret, err := taskService.Update(t)
			if err != nil {
				return fmt.Errorf("failed to update task %d, %v", id, err)
			}

			fmt.Printf("Task %d marked as 'ToDo'!\n", ret)
			return nil
		},
	}
)
