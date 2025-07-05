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
	"os"
	"path/filepath"
	"strings"

	"github.com/cgoesche/willdo/app"
	"github.com/cgoesche/willdo/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	conf       config.Config

	EditCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit a task or category",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	EditCmd.Flags().StringVar(&configFile, "config", "", "configuration file location")

	EditCmd.AddCommand(categoryCmd)
	EditCmd.AddCommand(taskCmd)
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
