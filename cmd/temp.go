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

package cmd

import (
	"fmt"

	"github.com/darox/sunly/internal/printer"
	"github.com/darox/sunly/pkg/swissmeteo"
	"github.com/darox/sunly/pkg/swisspost"
	"github.com/spf13/cobra"
)

// tempCmd represents the temp command.
var tempCmd = &cobra.Command{
	Use:   "temp",
	Short: "Returns the temperature of a location by providing a postal code",
	Long:  `Returns the temperature of a location by providing a postal code`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case rootCmd.PersistentFlags().Lookup("zip") != nil:
			l := rootCmd.PersistentFlags().Lookup("zip").Value.String()
			getCurrentTemperature(l)
		case rootCmd.PersistentFlags().Lookup("city") != nil:
			l := rootCmd.PersistentFlags().Lookup("city").Value.String()
			getCurrentTemperature(l)
		default:
			fmt.Println("Please provide a postal code or a city name")
		}
	},
}

func init() {
	rootCmd.AddCommand(tempCmd)
}

func getCurrentTemperature(locationIdentifier string) {
	// Create a new weather struct
	w := swissmeteo.Weather{}

	// Get the current weather
	c, err := w.GetCurrentWeather(locationIdentifier)
	if err != nil {
		fmt.Printf("Something went wrong when fetching the temperature: %s\n", err)
		return
	}

	// Create a new location data struct
	ld := swisspost.LocationData{}

	// Get the location data by zip code
	err = ld.GetLocationDataByZip(locationIdentifier)
	if err != nil {
		fmt.Printf("Something went wrong when fetching the location: %s\n", err)
		return
	}
	// Check if the zip code is valid
	if !ld.IsZipValid(zip) {
		fmt.Printf("The zip code %s is not valid\n", zip)
		return
	}

	// Extract the city name from the location data
	aCity := ld.Records[0].Fields.Ortbez18

	// Print the current temperature
	printer.PrintCurrentTemperature(c, zip, aCity)
}
