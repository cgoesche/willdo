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
	"github.com/cgoesche/willdo/cmd/add"
	"github.com/cgoesche/willdo/cmd/del"
	"github.com/cgoesche/willdo/cmd/edit"
	"github.com/cgoesche/willdo/cmd/list"
	"github.com/cgoesche/willdo/internal/bubbletea"
	"github.com/cgoesche/willdo/internal/config"
	"github.com/cgoesche/willdo/internal/database"

	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile   string
	databaseFile string
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
			client := database.NewClient()
			err := client.InitDB(databaseFile)
			if err != nil {
				return err
			}

			bubbletea.Run(client)
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

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Configuration file location")
	rootCmd.PersistentFlags().StringVar(&databaseFile, "database", config.SetDefault().Database.Path, "Database file path")

	viper.BindPFlag("database.path", rootCmd.PersistentFlags().Lookup("database"))

	rootCmd.AddCommand(add.AddCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(del.DeleteCmd)
	rootCmd.AddCommand(edit.EditCmd)
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(startCmd)
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
