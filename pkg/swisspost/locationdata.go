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

package swisspost

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// This package is using the official swiss post API to convert a zip code to a location and vice versa.

const (
	domain   = "https://swisspost.opendatasoft.com"
	pathZips = "/api/records/1.0/search/?dataset=plz_verzeichnis_v2&q=&rows=20&refine.gplz=%s"
	pathLoc  = "/api/records/1.0/search/?dataset=plz_verzeichnis_v2&q=&rows=20&refine.ortbez18=%s"
)

func (l *LocationData) GetLocDatByZip(zip string) (err error) {
	// Get the location data from the API

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(domain+pathZips, zip), nil)

	if err != nil {
		err = fmt.Errorf("error fetching location: %w", err)
		return err
	}

	// Execute the request
	c := http.DefaultClient
	resp, err := c.Do(req)

	// Check for errors
	if err != nil {
		return fmt.Errorf("error getting location data from API: %w", err)
	}

	// Close the body when we're done with it
	defer resp.Body.Close()

	// Decode the JSON respons
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		return err
	}
	// Return nil if everything went fine
	return nil
}

func (l *LocationData) GetLocationDataByName(name string) (err error) {
	// Get the location data from the API

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(domain+pathLoc, name), nil)

	if err != nil {
		err = fmt.Errorf("error fetching location: %w", err)
		return err
	}

	// Execute the request
	c := http.DefaultClient
	resp, err := c.Do(req)

	// Check for errors
	if err != nil {
		return fmt.Errorf("error getting location data from API: %w", err)
	}

	// Close the body when we're done with it
	defer resp.Body.Close()

	// Decode the JSON respons
	err = json.NewDecoder(resp.Body).Decode(&l)
	if err != nil {
		return err
	}
	// Return nil if everything went fine
	return nil
}

func (l *LocationData) ConvertZipToName(zip string) (name string, err error) {
	// Get the location data from the API
	err = l.GetLocDatByZip(zip)
	if err != nil {
		return name, err
	}

	// Return the location
	return l.Records[0].Fields.Ortbez18, nil
}

func (l *LocationData) ConvertCityToZips(name string) (zip []string, err error) {
	// Get the location data from the API
	err = l.GetLocationDataByName(name)
	if err != nil {
		return zip, err
	}

	for _, v := range l.Records {
		zip = append(zip, v.Fields.Postleitzahl)
	}

	// Return the zip codes
	return zip, nil
}

// Checks if the zip code is valid.
func (l *LocationData) IsZipValid(zip string) (valid bool) {
	if len(zip) != 4 {
		return false
	}

	if l.Nhits == 0 {
		return false
	}

	return true
}

func (l *LocationData) ExpandZipRange(zip string) (zips []string, err error) {
	// Get the location data from the API
	err = l.GetLocDatByZip(zip)
	if err != nil {
		err = fmt.Errorf("error getting location data from API: %w", err)
		return zips, err
	}

	for _, v := range l.Records {
		zips = append(zips, v.Fields.Postleitzahl)
	}
	// Return the zip codes
	return zips, nil
}

type LocationData struct {
	Nhits      int `json:"nhits"`
	Parameters struct {
		Dataset  string   `json:"dataset"`
		Q        string   `json:"q"`
		Rows     int      `json:"rows"`
		Start    int      `json:"start"`
		Facet    []string `json:"facet"`
		Format   string   `json:"format"`
		Timezone string   `json:"timezone"`
	} `json:"parameters"`
	Records []struct {
		Datasetid string `json:"datasetid"`
		Recordid  string `json:"recordid"`
		Fields    struct {
			Ortbez27     string    `json:"ortbez27"`
			GeoPoint2D   []float64 `json:"geo_point_2d"`
			PlzCoff      string    `json:"plz_coff"`
			RecArt       string    `json:"rec_art"`
			Sprachcode   int       `json:"sprachcode"`
			Bfsnr        int       `json:"bfsnr"`
			Kanton       string    `json:"kanton"`
			GiltAbDat    string    `json:"gilt_ab_dat"`
			Onrp         int       `json:"onrp"`
			Postleitzahl string    `json:"postleitzahl"`
			Gplz         int       `json:"gplz"`
			PlzBriefzust int       `json:"plz_briefzust"`
			Ortbez18     string    `json:"ortbez18"`
			BriefzDurch  int       `json:"briefz_durch"`
			PlzZz        string    `json:"plz_zz"`
			GeoShape     struct {
				Coordinates [][][]float64 `json:"coordinates"`
				Type        string        `json:"type"`
			} `json:"geo_shape"`
			PlzTyp int `json:"plz_typ"`
		} `json:"fields"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		RecordTimestamp time.Time `json:"record_timestamp"`
	} `json:"records"`
	FacetGroups []struct {
		Name   string `json:"name"`
		Facets []struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
			State string `json:"state"`
			Path  string `json:"path"`
		} `json:"facets"`
	} `json:"facet_groups"`
}
