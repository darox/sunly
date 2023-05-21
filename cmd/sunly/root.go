/*
Sunly
Copyright (C) 2023 Dario Mader

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

package sunly

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var (
	rootCmd = &cobra.Command{
		Use:   "sunly",
		Short: "sunly - a CLI tool to get the weather from MeteoSwiss",
		Long: `sunly - a CLI tool to get the weather from MeteoSwiss. Examples:
		sunly get current --zip 8000
		sunly get forcecast --city Zurich`,
	}
	pZip  string
	pCity string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sunly.yaml)")

	rootCmd.PersistentFlags().StringVar(&pZip, "zip", "", "Postal code of the location")
	rootCmd.PersistentFlags().StringVar(&pCity, "city", "", "City name")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
