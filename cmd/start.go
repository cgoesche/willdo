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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	startCmd = &cobra.Command{
		Use:   "start [ID]",
		Short: "Start a task",
		Long: `There is not much more to say about this or 
are you looking for the entire commit history ?`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := viper.GetViper().GetString("database.path")
			client := database.NewClient()

			err := client.InitDB(path)
			if err != nil {
				return err
			}

			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task ID")
			}

			_, err = client.StartTask(id)
			if err != nil {
				return fmt.Errorf("failed to start task, %v", err)
			}

			fmt.Printf("Task %d started!\n", id)
			return nil
		},
	}
)
