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
		case rootCmd.PersistentFlags().Lookup("zip").Value.String() != "":
			zip := rootCmd.PersistentFlags().Lookup("zip").Value.String()
			err := getCurrentTemperature("", zip)
			if err != nil {
				fmt.Printf("Something went wrong when fetching the temperature: %s\n", err)
			}
		case rootCmd.PersistentFlags().Lookup("city").Value.String() != "":
			city := rootCmd.PersistentFlags().Lookup("city").Value.String()
			err := getCurrentTemperature(city, "")
			if err != nil {
				fmt.Printf("Something went wrong when fetching the temperature: %s\n", err)
			}
		default:
			fmt.Println("Please provide a postal code or a city name")
		}
	},
}

func init() {
	rootCmd.AddCommand(tempCmd)
}

func getCurrentTemperature(city string, zip string) error {
	// Create a new location data struct
	d := swisspost.LocationData{}

	// Create a new weather struct
	w := swissmeteo.Weather{}
	zips := []string{}

	if city != "" {
		var err error
		zips, err = d.ConvertCityToZips(city)
		if err != nil {
			fmt.Printf("Something went wrong when fetching the location: %s\n", err)
			return err
		}
	} else {
		var err error
		zips, err = d.ExpandZipRange(zip)
		fmt.Println(zips)
		if err != nil {

		}
	}

	var c swissmeteo.CurrentWeather

	c, matchedZip, err := w.GetCurrentWeather(zips)
	if err != nil {
		return err
	}

	// Get the location data by zip code
	err = d.GetLocDatByZip(matchedZip)
	if err != nil {
		fmt.Printf("Something went wrong when fetching the location: %s\n", err)
		return err
	}
	// Check if the zip code is valid
	if !d.IsZipValid(matchedZip) {
		fmt.Printf("The zip code %s is not valid\n", zip)
		return nil
	}

	// Extract the city name from the location data
	aCity := d.Records[0].Fields.Ortbez18

	// Print the current temperature
	printer.PrintCurrentTemperature(c, matchedZip, aCity)

	return nil
}
