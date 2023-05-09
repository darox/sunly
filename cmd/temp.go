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
	"time"

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
			getCurrentTemperature(zip)
		//TODO: Add location flag
		default:
			fmt.Println("Please provide a zip code")
		}
	},
}

func init() {
	rootCmd.AddCommand(tempCmd)
}

func getCurrentTemperature(zip string) {
	// Create a new weather object
	w := swissmeteo.Weather{}

	// Get the current temperature for the given zip code
	temperature, u, err := w.GetCurrentTemperature(zip)
	if err != nil {
		fmt.Printf("Something went wrong when fetching the temperature: %s\n", err)
		return
	}

	// Convert time to a human readable format
	h := time.Unix(u/1000, 0)
	updatedAt := h.Format("15:04 02.01.2006")

	// Create a new location object
	ld := swisspost.LocationData{}

	// Get the location data
	err = ld.GetLocationDataByZip(zip)

	if err != nil {
		fmt.Printf("Something went wrong when fetching the location: %s\n", err)
		return
	}

	locationName := ld.Records[0].Fields.Ortbez18

	// Check if the zip code is valid

	if !ld.IsZipValid(zip) {
		fmt.Printf("The zip code %s is not valid\n", zip)
		return
	}

	printer.PrintCurrentTemperature(zip, locationName, temperature, updatedAt)
}
