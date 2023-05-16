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

package swissmeteo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	apiURL = "https://app-prod-ws.meteoswiss-app.ch/v1/plzDetail?plz=%s"
)

// Gets the weather data from the API and decodes it into the Weather struct.
func (w *Weather) getWeatherData(zip string) error {
	// The meteoswiss API only accepts a zip code with a tailing 00
	z := fmt.Sprintf("%s00", zip)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// Create a new request with the context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(apiURL, z), nil)

	if err != nil {
		err = fmt.Errorf("error fetching weather data: %w", err)
		return err
	}

	// Execute the request
	c := http.DefaultClient
	resp, err := c.Do(req)

	// Check for errors
	if err != nil {
		return fmt.Errorf("error getting weather data from API: %w", err)
	}

	// Close the body when we're done with it
	defer resp.Body.Close()

	// Decode the JSON respons
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return err
	}

	return nil
}

func (w *Weather) GetCurrentWeather(zip []string) (c CurrentWeather, matchedZip string, err error) {
	for _, z := range zip {
		err := w.getWeatherData(z)
		if err != nil {
			err = fmt.Errorf("error getting weather data: %w", err)
			return w.CurrentWeather, "", err
		}
		if w.CurrentWeather.Time != 1684587000000 {
			return w.CurrentWeather, z, nil
		}
		return w.CurrentWeather, "", nil
	}
	return w.CurrentWeather, "", nil
}

type CurrentWeather struct {
	Time        int64   `json:"time"`
	Icon        int     `json:"icon"`
	IconV2      int     `json:"iconV2"`
	Temperature float64 `json:"temperature"`
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
