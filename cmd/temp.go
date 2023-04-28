/*
Copyright © 2023 Dario Mader maderdario@gmail.com
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

const (
	API_URL = "https://app-prod-ws.meteoswiss-app.ch/v1/plzDetail?plz=%s"
)

// tempCmd represents the temp command
var tempCmd = &cobra.Command{
	Use:   "temp",
	Short: "Returns the temperature of a location by providing a postal code",
	Long:  `Returns the temperature of a location by providing a postal code`,
	Run: func(cmd *cobra.Command, args []string) {
		getTemp(args[0])
	},
}

func init() {
	rootCmd.AddCommand(tempCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tempCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tempCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getTemp(zip string) {
	// Fetch the weather from the API
	z := zip + "00"
	w, err := http.Get(fmt.Sprintf(API_URL, z))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer w.Body.Close()
	var weather Weather
	// Decode the JSON response into our struct type.
	err = json.NewDecoder(w.Body).Decode(&weather)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print the temperature
	fmt.Printf("The temperature in %s is %0.1f°C", zip, weather.CurrentWeather.Temperature)
}

// Weather struct to hold the JSON response
type Weather struct {
	CurrentWeather struct {
		Time        int64   `json:"time"`
		Icon        int     `json:"icon"`
		IconV2      int     `json:"iconV2"`
		Temperature float64 `json:"temperature"`
	} `json:"currentWeather"`
	Forecast []struct {
		DayDate        string  `json:"dayDate"`
		IconDay        int     `json:"iconDay"`
		IconDayV2      int     `json:"iconDayV2"`
		TemperatureMax int     `json:"temperatureMax"`
		TemperatureMin int     `json:"temperatureMin"`
		Precipitation  float64 `json:"precipitation"`
	} `json:"forecast"`
	Warnings         []any `json:"warnings"`
	WarningsOverview []any `json:"warningsOverview"`
	Graph            struct {
		Start               int64     `json:"start"`
		StartLowResolution  int64     `json:"startLowResolution"`
		Precipitation10M    []float64 `json:"precipitation10m"`
		PrecipitationMin10M []float64 `json:"precipitationMin10m"`
		PrecipitationMax10M []float64 `json:"precipitationMax10m"`
		WeatherIcon3H       []int     `json:"weatherIcon3h"`
		WeatherIcon3HV2     []int     `json:"weatherIcon3hV2"`
		WindDirection3H     []int     `json:"windDirection3h"`
		WindSpeed3H         []float64 `json:"windSpeed3h"`
		Sunrise             []int64   `json:"sunrise"`
		Sunset              []int64   `json:"sunset"`
		TemperatureMin1H    []float64 `json:"temperatureMin1h"`
		TemperatureMax1H    []float64 `json:"temperatureMax1h"`
		TemperatureMean1H   []float64 `json:"temperatureMean1h"`
		Precipitation1H     []float64 `json:"precipitation1h"`
		PrecipitationMin1H  []float64 `json:"precipitationMin1h"`
		PrecipitationMax1H  []float64 `json:"precipitationMax1h"`
	} `json:"graph"`
}
