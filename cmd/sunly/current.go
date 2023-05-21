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
	"fmt"

	printer "github.com/darox/sunly/internal/printer"
	"github.com/darox/sunly/pkg/swissmeteo"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Returns the current weather of a location by providing a zip code or a city name",
	Long: `Returns the current weather of a location by providing a zip code or a city name. Examples:
	sunly current --zip 8000
	sunly current --city Zurich`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case rootCmd.PersistentFlags().Lookup("zip").Value.String() != "":
			zip := rootCmd.PersistentFlags().Lookup("zip").Value.String()
			err := GetCurrentWeather("", zip)
			if err != nil {
				fmt.Printf("Something went wrong when fetching the current weather: %s\n", err)
			}
		case rootCmd.PersistentFlags().Lookup("city").Value.String() != "":
			city := rootCmd.PersistentFlags().Lookup("city").Value.String()
			err := GetCurrentWeather(city, "")
			if err != nil {
				fmt.Printf("Something went wrong when fetching the current weather: %s\n", err)
			}
		default:
			fmt.Println("Please provide a postal code or a city name")
		}
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

func GetCurrentWeather(city string, zip string) error {

	// Creat a new current weather struct
	w := swissmeteo.NewCurrentWeather(zip)

	// Get the current weather data
	err := w.GetCurrentWeather()
	if err != nil {
		return err
	}
	/*
		if err != nil {
			s, err := swissmeteo.GetSupportedZipCodes(zip)
			if err != nil {
				fmt.Printf("Something went wrong when fetching the current weather: %s\n", err)
				return err
			}
			fmt.Printf(`Something went wrong when fetching the current weather: %s.
			The MeteoSwiss API also doesn't return data for every zip code. Try one of the following: %s\n
			`, err, s)
			return err
		} */

	printer.PrintCurrentWeather(*w)

	return nil
}
