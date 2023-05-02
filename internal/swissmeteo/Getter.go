package swissmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	API_URL = "https://app-prod-ws.meteoswiss-app.ch/v1/plzDetail?plz=%s"
)

func getWeather(zip string) (w Weather, err error) {
	// The meteoswiss API only accepts a zip code with a tailing 00
	z := zip + "00"
	// Fetch the weather from the API
	r, err := http.Get(fmt.Sprintf(API_URL, z))
	// Check for errors
	if err != nil {
		err = fmt.Errorf("Error fetching weather: %s", err)
		return w, err
	}
	// Close the body when we're done with it
	defer r.Body.Close()

	// Decode the JSON respons
	err = json.NewDecoder(r.Body).Decode(&w)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Return the weather
	return w, nil
}

func GetTemp(zip string) (t float64, u int64, err error) {
	// Get the weather
	w, err := getWeather(zip)
	if err != nil {
		return t, u, err
	}
	// Return the temperature
	return w.CurrentWeather.Temperature, w.CurrentWeather.Time, nil
}

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
