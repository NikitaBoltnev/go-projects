// Package weather provides functionality to retrieve weather forecasts
// from the wttr.in service based on geographic location and output format.
package weather

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"weather-cli/geo"
)

var ErrWrongFormat = errors.New("WRONG_FORMAT")

// GetWeather fetches weather for a city in the specified format (1-4).
func GetWeather(geo geo.GeoData, format int) (string, error) {
	baseUrl, err := url.Parse("https://wttr.in/" + geo.City)
	if err != nil {
		return "", errors.New("ERROR_URL")
	}

	if format < 1 || format > 4 {
		return "", ErrWrongFormat
	}

	params := url.Values{}
	params.Add("format", fmt.Sprint(format))
	baseUrl.RawQuery = params.Encode()
	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return "", errors.New("ERROR_HTTP")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("ERROR_BODY")
	}

	return string(body), nil
}
