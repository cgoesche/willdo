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
package config

import (
	"os"
	"path/filepath"
)

type Database struct {
	Type     string `mapstructure:"type"`
	NetAddr  string `mapstructure:"address"`
	Port     int    `mapstructure:"port"`
	Filepath string `mapstructure:"filepath"`
}

type Config struct {
	Database Database `mapstructure:"database"`
}

func SetDefault() Config {
	dir, _ := os.UserHomeDir()
	dbFilePath := filepath.Join(dir, "willdo.db")

	return Config{
		Database: Database{
			Type:     "sqlite",
			Filepath: dbFilePath,
		},
	}
}
