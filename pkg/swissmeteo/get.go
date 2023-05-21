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
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL              = "https://www.meteoschweiz.admin.ch/product/output/"
	versionsPath         = "versions.json"
	forecastPath         = "weather-widget/forecast/version__%s/en/%s00.json"
	supportedZipCodesURL = "https://www.meteoschweiz.admin.ch/static/product/resources/local-forecast-search/%s.json"
)

func NewWeatherVersions() *WeatherVersions {
	return &WeatherVersions{}
}

// Instantiate a new current weather struct
func NewCurrentWeather(zip string) *CurrentWeather {
	return &CurrentWeather{
		Zip: zip,
	}
}

// Instantiate a new forecast weather struct
func NewForecastWeather(zip string) *ForecastWeather {
	return &ForecastWeather{
		Zip: zip,
	}
}

// Instantiate a new weather struct
func NewWeather(version string, zip string) *Weather {
	return &Weather{
		Version: version,
		Zip:     zip,
	}
}

// Fetch the weather data from the given URL
func fetchDataFromApi(u url.URL) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	c := http.DefaultClient
	resp, err := c.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading HTTP response: %w", err)
	}

	return body, nil

}

// Get the weather versions by calling the fetchDataFromApi function
func (w *WeatherVersions) GetWeatherVersions() error {

	baseURL, err := url.Parse(baseURL)

	if err != nil {
		return fmt.Errorf("error parsing base URL: %w", err)
	}

	relativeURL, err := url.Parse(versionsPath)

	if err != nil {
		return fmt.Errorf("error parsing relative URL: %w", err)
	}

	newURL := baseURL.ResolveReference(relativeURL)

	resp, err := fetchDataFromApi(*newURL)

	if err != nil {
		return fmt.Errorf("error fetching weather versions: %w", err)
	}

	err = json.Unmarshal(resp, &w)

	if err != nil {
		return fmt.Errorf("error unmarshalling weather versions: %w", err)
	}

	return nil
}

// Get the weather data by calling the fetchDataFromApi function
func (w *Weather) GetWeather(version string, zip string) error {

	baseURL, err := url.Parse(baseURL)

	if err != nil {
		return fmt.Errorf("error parsing base URL: %w", err)
	}

	relativeURL := fmt.Sprintf(forecastPath, version, zip)

	path, err := url.Parse(relativeURL)

	if err != nil {
		return fmt.Errorf("error parsing relative URL: %w", err)
	}

	newURL := baseURL.ResolveReference(path)

	resp, err := fetchDataFromApi(*newURL)

	if err != nil {
		return fmt.Errorf("error fetching weather data: %w", err)
	}

	err = json.Unmarshal(resp, &w)

	if err != nil {
		return fmt.Errorf("error unmarshalling weather data: %w", err)
	}

	return nil
}

// Get the current weather by calling the GetWeatherData function
func (cw *CurrentWeather) GetCurrentWeather() error {

	wv := NewWeatherVersions()

	wv.GetWeatherVersions()

	w := NewWeather(wv.WeatherWidgetForecast, cw.Zip)

	w.GetWeather(wv.WeatherWidgetForecast, cw.Zip)

	cw.Current = w.Data.Current
	cw.CityName = w.Data.CityName
	cw.Timestamp = w.Data.Timestamp

	return nil
}

func GetSupportedZipCodes(zip string) ([]map[string]string, error) {
	// Strip the last two digits from the zip code
	z := zip[:len(zip)-2]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(supportedZipCodesURL, z), nil)

	if err != nil {
		err = fmt.Errorf("error fetching weather data: %w", err)
		return nil, err
	}

	c := http.DefaultClient
	resp, err := c.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error getting weather data from API: %w", err)
	}

	defer resp.Body.Close()

	var s SupportedZipCodes

	err = json.NewDecoder(resp.Body).Decode(&s)

	if err != nil {
		return nil, err
	}

	supportedZips := make([]map[string]string, len(s))
	for _, v := range s {
		s := strings.Split(v, ";")
		supportedZips = append(supportedZips, map[string]string{
			s[7]: s[0]})
	}

	return supportedZips, nil
}

type SupportedZipCodes []string

type CurrentWeather struct {
	Current struct {
		Temperature     string `json:"temperature"`
		WeatherSymbolID string `json:"weather_symbol_id"`
	} `json:"current"`
	CityName  string `json:"city_name"`
	Timestamp int    `json:"timestamp"`
	Zip       string `json:"zip"`
}

type ForecastWeather struct {
	Forecast []struct {
		Noon            int64  `json:"noon"`
		TempHigh        string `json:"temp_high"`
		Weekday         string `json:"weekday"`
		TempLow         string `json:"temp_low"`
		WeatherSymbolID string `json:"weather_symbol_id"`
	} `json:"forecasts"`
	CityName  string `json:"city_name"`
	Timestamp int    `json:"timestamp"`
	Zip       string `json:"zip"`
}

type Weather struct {
	Data struct {
		Altitude int    `json:"altitude"`
		CityName string `json:"city_name"`
		Current  struct {
			Temperature     string `json:"temperature"`
			WeatherSymbolID string `json:"weather_symbol_id"`
		} `json:"current"`
		WeatherSymbolID int    `json:"weather_symbol_id"`
		LocationID      string `json:"location_id"`
		Forecasts       []struct {
			Noon            int64  `json:"noon"`
			TempHigh        string `json:"temp_high"`
			Weekday         string `json:"weekday"`
			TempLow         string `json:"temp_low"`
			WeatherSymbolID string `json:"weather_symbol_id"`
		} `json:"forecasts"`
		Timestamp int `json:"timestamp"`
	} `json:"data"`
	Config struct {
		Name      string `json:"name"`
		Language  string `json:"language"`
		Version   string `json:"version"`
		Timestamp int    `json:"timestamp"`
	} `json:"config"`
	Zip     string `json:"zip"`
	Version string `json:"version"`
}

type WeatherVersions struct {
	PollenAnimationAll                            string `json:"pollen/animation/all"`
	SnowMapsPrecipitation24H                      string `json:"snow/maps/precipitation-24h"`
	CosmoTemperatureAnimation                     string `json:"cosmo/temperature/animation"`
	SnowMapsPrecipitation12H                      string `json:"snow/maps/precipitation-12h"`
	TeaserImageCloudCoverMap                      string `json:"teaser-image/cloud-cover-map"`
	WeatherOutlookItWest                          string `json:"weather-outlook/it/west"`
	TeaserImageMeasuredValues                     string `json:"teaser-image/measured-values"`
	Danger                                        string `json:"danger"`
	WeatherOutlookDeSouth                         string `json:"weather-outlook/de/south"`
	WeatherReportFrSouth                          string `json:"weather-report/fr/south"`
	TeaserImageLakesAndAirfieldsMap               string `json:"teaser-image/lakes-and-airfields-map"`
	TeaserImageNaturalHazardsMap                  string `json:"teaser-image/natural-hazards-map"`
	ForecastOverviewChart                         string `json:"forecast-overview-chart"`
	SatelliteHrvImages                            string `json:"satellite/hrv/images"`
	ClimateAtmosphereRadioAnomPzintps0LastClimrep string `json:"climate-atmosphere-radio-anom/pzintps0/last/climrep"`
	TeaserImageSnowHeightMap                      string `json:"teaser-image/snow-height-map"`
	WeatherRegionOverview                         string `json:"weather-region-overview"`
	CosmoCloudCoverAnimation                      string `json:"cosmo/cloud-cover/animation"`
	WeatherOutlookFrNorth                         string `json:"weather-outlook/fr/north"`
	SnowMapsPrecipitation48H                      string `json:"snow/maps/precipitation-48h"`
	ClimateMonitorSunshine                        string `json:"climate-monitor/sunshine"`
	CosmoWind10MImages                            string `json:"cosmo/wind-10m/images"`
	Lightning                                     string `json:"lightning"`
	ClimateIndicatorsHEATDLUZ                     string `json:"climate-indicators/HEATD/LUZ"`
	PollenImagesBetu                              string `json:"pollen/images/betu"`
	WeatherWidgetForecast                         string `json:"weather-widget/forecast"`
	ForecastMap                                   string `json:"forecast-map"`
	TeaserImageSatelliteAnimationHrv              string `json:"teaser-image/satellite-animation-hrv"`
	TemperatureAnimation                          string `json:"temperature/animation"`
	TeaserImageAnimationWind10M                   string `json:"teaser-image/animation-wind-10m"`
	RoadConditionIt                               string `json:"road-condition/it"`
	WeatherReportDeNorth                          string `json:"weather-report/de/north"`
	ClimateIndicatorsTNNEU                        string `json:"climate-indicators/TN/NEU"`
	WeatherWidgetDanger                           string `json:"weather-widget/danger"`
	GeneralsituationMapFr                         string `json:"generalsituation/map/fr"`
	PollenAnimationCory                           string `json:"pollen/animation/cory"`
	CosmoWindGusts10MAnimation                    string `json:"cosmo/wind-gusts-10m/animation"`
	ClimatePhenoSeriesLongLZLIW                   string `json:"climate-pheno-series-long/LZLIW"`
	SnowImagePrecipitation72H                     string `json:"snow/image/precipitation-72h"`
	PollenTextDe                                  string `json:"pollen/text/de"`
	CosmoWindGustsAnimation                       string `json:"cosmo/wind/gusts-animation"`
	ClimateAtmosphereRadioAnomPzinfrs0LastClimrep string `json:"climate-atmosphere-radio-anom/pzinfrs0/last/climrep"`
	WeatherOutlookItNorth                         string `json:"weather-outlook/it/north"`
	TeaserImageClimateMonitorSunshineMap          string `json:"teaser-image/climate-monitor-sunshine-map"`
	GeneralsituationTextIt                        string `json:"generalsituation/text/it"`
	SatelliteCloudIceAnimation                    string `json:"satellite/cloud-ice/animation"`
	SatelliteCloudIceImages                       string `json:"satellite/cloud-ice/images"`
	ClimateMonitorDiagram                         string `json:"climate-monitor/diagram"`
	TeaserImageSnowForecastMap                    string `json:"teaser-image/snow-forecast-map"`
	CosmoWind2000MImages                          string `json:"cosmo/wind-2000m/images"`
	ForecastChart                                 string `json:"forecast-chart"`
	WeatherReportItNorth                          string `json:"weather-report/it/north"`
	SatelliteEuropeIrAnimation                    string `json:"satellite/europe-ir/animation"`
	ClimateEvolution                              string `json:"climate-evolution"`
	NaturalHazardBulletin                         string `json:"natural-hazard-bulletin"`
	RadioSoundingsEmagram                         string `json:"radio-soundings/emagram"`
	WeatherReportItSouth                          string `json:"weather-report/it/south"`
	WeatherPill                                   string `json:"weather-pill"`
	WeatherReportItWest                           string `json:"weather-report/it/west"`
	ForecastText                                  string `json:"forecast-text"`
	RadioSoundingsDecoded                         string `json:"radio-soundings/decoded"`
	CosmoCloudCoverForecast                       string `json:"cosmo/cloud-cover/forecast"`
	GeneralsituationMapDe                         string `json:"generalsituation/map/de"`
	SnowMapsHeights                               string `json:"snow/maps/heights"`
	SatelliteCloudCover                           string `json:"satellite/cloud-cover"`
	ClimateIndicatorsCDDOTL                       string `json:"climate-indicators/CDD/OTL"`
	PollenAnimationBetu                           string `json:"pollen/animation/betu"`
	ClimateMonitorPrecipitation                   string `json:"climate-monitor/precipitation"`
	UvIndex                                       string `json:"uv-index"`
	ClimateIndicatorsCDDSTG                       string `json:"climate-indicators/CDD/STG"`
	SatelliteWorldMosaicAnimation                 string `json:"satellite/world-mosaic/animation"`
	TeaserImageSatelliteAnimationEuropeIr         string `json:"teaser-image/satellite-animation-europe-ir"`
	ClimateIndicatorsHDD20GVE                     string `json:"climate-indicators/HDD20/GVE"`
	TeaserImageTemperatureMap                     string `json:"teaser-image/temperature-map"`
	ClimateIndicatorsSDEIN                        string `json:"climate-indicators/SD/EIN"`
	SatelliteHrvAnimation                         string `json:"satellite/hrv/animation"`
	PollenImagesCory                              string `json:"pollen/images/cory"`
	ClimateIndicatorsCDDBAS                       string `json:"climate-indicators/CDD/BAS"`
	TeaserImageSatelliteAnimationCloudIce         string `json:"teaser-image/satellite-animation-cloud-ice"`
	CosmoWind2000MForecast                        string `json:"cosmo/wind/2000m/forecast"`
	WeatherReportFrWest                           string `json:"weather-report/fr/west"`
	RoadConditionDe                               string `json:"road-condition/de"`
	CosmoWind10MForecast                          string `json:"cosmo/wind/10m/forecast"`
	WeatherOutlookDeNorth                         string `json:"weather-outlook/de/north"`
	IncaPrecipitationRate                         string `json:"inca/precipitation/rate"`
	PrecipitationAnimation                        string `json:"precipitation/animation"`
	CosmoWind2000MAnimation                       string `json:"cosmo/wind-2000m/animation"`
	CosmoWindGusts10MImages                       string `json:"cosmo/wind-gusts-10m/images"`
	WeatherOutlookDeWest                          string `json:"weather-outlook/de/west"`
	SnowMapsPrecipitation72H                      string `json:"snow/maps/precipitation-72h"`
	WeatherReportDeSouth                          string `json:"weather-report/de/south"`
	TeaserImageClimateMonitorTemperatureMap       string `json:"teaser-image/climate-monitor-temperature-map"`
	IncaPrecipitationTypeFreezingRain             string `json:"inca/precipitation/type/freezing-rain"`
	CosmoWind2000MAnimation0                      string `json:"cosmo/wind/2000m-animation"`
	SaharanDustEvents                             string `json:"saharan-dust/events"`
	GeneralsituationMapIt                         string `json:"generalsituation/map/it"`
	Sov                                           string `json:"sov"`
	TeaserImageSevereWeatherMap                   string `json:"teaser-image/severe-weather-map"`
	PollenTextIt                                  string `json:"pollen/text/it"`
	TeaserImageAnimationWind2000M                 string `json:"teaser-image/animation-wind-2000m"`
	CosmoTemperatureImages                        string `json:"cosmo/temperature/images"`
	PollenTabs                                    string `json:"pollen/tabs"`
	ClimateMonitorTemperature                     string `json:"climate-monitor/temperature"`
	PollenImagesAmbr                              string `json:"pollen/images/ambr"`
	IncaPrecipitationTypeSnowrain                 string `json:"inca/precipitation/type/snowrain"`
	TeaserImagePollenMap                          string `json:"teaser-image/pollen-map"`
	PollenAnimationAlnu                           string `json:"pollen/animation/alnu"`
	CosmoPrecipitationSum24HPossib                string `json:"cosmo/precipitation-sum-24h-possib"`
	SnowImageHeights                              string `json:"snow/image/heights"`
	WeatherOutlookFrSouth                         string `json:"weather-outlook/fr/south"`
	ClimateIndicatorsFD0DAV                       string `json:"climate-indicators/FD0/DAV"`
	ClimateIndicatorsID0JUN                       string `json:"climate-indicators/ID0/JUN"`
	ClimatePhenoSeriesSpringindex                 string `json:"climate-pheno-series-springindex"`
	WeatherOverview                               string `json:"weather-overview"`
	DataAvailability                              string `json:"data-availability"`
	IncaPrecipitationTypeSnow                     string `json:"inca/precipitation/type/snow"`
	GeneralsituationTextDe                        string `json:"generalsituation/text/de"`
	SatelliteEuropeIrImages                       string `json:"satellite/europe-ir/images"`
	ClimatePhenoSeriesLongPGE                     string `json:"climate-pheno-series-long/PGE"`
	CosmoWindGustsForecast                        string `json:"cosmo/wind/gusts/forecast"`
	WeatherReportFrNorth                          string `json:"weather-report/fr/north"`
	ClimateIndicatorsCDDGVE                       string `json:"climate-indicators/CDD/GVE"`
	TeaserImagePrecipitationMap                   string `json:"teaser-image/precipitation-map"`
	WeatherOutlookItSouth                         string `json:"weather-outlook/it/south"`
	AltitudeLevels                                string `json:"altitude-levels"`
	PollenAnimationPoac                           string `json:"pollen/animation/poac"`
	PollenImagesPoac                              string `json:"pollen/images/poac"`
	WeatherReportDeWest                           string `json:"weather-report/de/west"`
	CosmoWind10MAnimation                         string `json:"cosmo/wind-10m/animation"`
	RoadConditionFr                               string `json:"road-condition/fr"`
	TeaserImageUvIndexMap                         string `json:"teaser-image/uv-index-map"`
	SatelliteWorldMosaicImages                    string `json:"satellite/world-mosaic/images"`
	ClimateSeriesCitylandtemp                     string `json:"climate-series-citylandtemp"`
	PollenImagesAll                               string `json:"pollen/images/all"`
	PollenAnimationAmbr                           string `json:"pollen/animation/ambr"`
	WeatherOutlookFrWest                          string `json:"weather-outlook/fr/west"`
	GeneralsituationTextFr                        string `json:"generalsituation/text/fr"`
	PollenImagesAlnu                              string `json:"pollen/images/alnu"`
	CosmoWind10MAnimation0                        string `json:"cosmo/wind/10m-animation"`
	TeaserImageSatelliteAnimationWorldMosaic      string `json:"teaser-image/satellite-animation-world-mosaic"`
	IncaTemperature                               string `json:"inca/temperature"`
	SnowImagePrecipitation12H                     string `json:"snow/image/precipitation-12h"`
	CosmoPrecipitationSum48H                      string `json:"cosmo/precipitation-sum-48h"`
	SnowImagePrecipitation24H                     string `json:"snow/image/precipitation-24h"`
	TeaserImageAnimationWindGusts                 string `json:"teaser-image/animation-wind-gusts"`
	PollenTextFr                                  string `json:"pollen/text/fr"`
	TeaserImageClimateMonitorPrecipitationMap     string `json:"teaser-image/climate-monitor-precipitation-map"`
	CosmoPrecipitationSum24H                      string `json:"cosmo/precipitation-sum-24h"`
	SnowImagePrecipitation48H                     string `json:"snow/image/precipitation-48h"`
}
