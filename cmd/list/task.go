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
package list

import (
	"fmt"

	"github.com/cgoesche/willdo/internal/bubbletea"
	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/models"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	categoryID   int64
	taskID       int64
	listAllTasks bool

	taskCmd = &cobra.Command{
		Use:   "task",
		Short: "List tasks",
		Long: `There is not much more to say about this or 
are you looking for the entire commit history ?`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := viper.GetViper().GetString("database.path")
			client := database.NewClient()
			err := client.InitDB(path)
			if err != nil {
				return fmt.Errorf("%v %v", database.ErrOpenDatabase, err)
			}

			if listAllTasks {
				bubbletea.Run(client)
			}

			if taskID > 0 {
				if err := listTaskFromID(client, taskID); err != nil {
					return fmt.Errorf("cannot list task with ID %d, %v", taskID, err)
				}
			}
			return nil
		},
	}
)

func init() {
	taskCmd.Flags().Int64VarP(&taskID, "id", "i", 0, "List a specific task")
	taskCmd.Flags().Int64VarP(&categoryID, "category", "c", 0, "List tasks from a specific category")
	taskCmd.Flags().BoolVarP(&listAllTasks, "all", "a", false, "List tasks from all categories")

	taskCmd.MarkFlagsOneRequired("id", "category", "all")
	taskCmd.MarkFlagsMutuallyExclusive("id", "category", "all")
}

func listTaskFromID(c *database.Client, id int64) error {
	var t models.Task

	t, err := c.QueryTaskFromID(id)
	if err != nil {
		return err
	}

	var style = lipgloss.NewStyle().
		Background(lipgloss.Color("42")).
		Foreground(lipgloss.Color("170")).
		Padding(0, 1)

	fmt.Println(style.Render("Hello, kitty"))

	fmt.Printf("\nID\tTitle\n")
	fmt.Printf("-------------------------\n")
	fmt.Printf("%d\t%s\n", t.ID, t.Title)

	return nil
}
