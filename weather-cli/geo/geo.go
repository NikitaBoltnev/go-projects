// Package geo provides location lookup by city name or IP address.
package geo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type GeoData struct {
	City string `json:"city"`
}

type CityPopulationResponce struct {
	Error bool `json:"error"`
}

var ErrNoCity = errors.New("NOCITY")
var ErrNot200 = errors.New("NOT200")

// GetMyLocation returns the user's location.
// If city is provided and valid, uses it. Otherwise, detects via IP.
func GetMyLocation(city string) (*GeoData, error) {
	if city != "" {
		isCity := checkCity(city)
		if !isCity {
			return nil, ErrNoCity
		}
		return &GeoData{
			City: city,
		}, nil
	}
	resp, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		return nil, ErrNot200
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geo GeoData
	json.Unmarshal(body, &geo)

	return &geo, nil
}

// checkCity validates if the city exists using an external API.
func checkCity(city string) bool {
	postBody, _ := json.Marshal(map[string]string{
		"city": city,
	})

	resp, err := http.Post("https://countriesnow.space/api/v0.1/countries/population/cities",
		"application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var populationResponce CityPopulationResponce
	json.Unmarshal(body, &populationResponce)

	return !populationResponce.Error
}
