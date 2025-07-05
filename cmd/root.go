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
	"path/filepath"
	"strings"

	"github.com/cgoesche/willdo/app"
	"github.com/cgoesche/willdo/internal/bubbletea"
	"github.com/cgoesche/willdo/internal/config"
	"github.com/cgoesche/willdo/internal/database"
	"github.com/cgoesche/willdo/internal/modules/category"
	"github.com/cgoesche/willdo/internal/modules/task"

	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile   string
	taskID       int64
	categoryID   int64
	categoryName string
	showAllTasks bool
	conf         config.Config

	rootCmd = &cobra.Command{
		Use:     "willdo",
		Short:   "A featureful command line to-do list manager",
		Long:    `willdo is a featureful command line to-do list manager.`,
		Version: app.Version,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
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
			} else {
				showAllTasks = true
			}

			m := bubbletea.InitialModel()
			m.TaskService = taskService
			m.CategoryService = catService
			m.CatNameToIDMap = category.NewCategoryNameToIDMap(cats)
			m.CatIDToNameMap = category.NewCategoryIDToNameMap(cats)
			m.Categories = cats
			m.ShowAllTasks = showAllTasks
			m.SelectedCategory = categoryID

			bubbletea.Run(m)

			return nil
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&showAllTasks, "all", "a", false, "Show tasks from all categories")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Configuration file location")
	rootCmd.PersistentFlags().StringVarP(&categoryName, "category", "c", "", "Category to list tasks of")
	rootCmd.MarkFlagsMutuallyExclusive("all", "category")

	rootCmd.AddCommand(categoryCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		configDir, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		var configPath = filepath.Join(configDir, app.Name)
		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
	}

	viper.SetEnvPrefix(app.Name)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
		os.Exit(1)
	}
}
