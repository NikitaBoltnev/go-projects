// Package weather_test contains tests for the weather package.
package weather_test

import (
	"strings"
	"testing"
	"weather-cli/geo"
	"weather-cli/weather"
)

// TestGetWeather verifies that weather data is fetched for a valid city and format.
func TestGetWeather(t *testing.T) {
	expected := "London"

	geo := geo.GeoData{
		City: expected,
	}
	format := 3

	got, err := weather.GetWeather(geo, format)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !strings.Contains(got, expected) {
		t.Errorf("Expected response to contain %q, got: %q", expected, got)
	}
}

// Test cases for invalid format values.
var testCases = []struct {
	name   string
	format int
}{
	{name: "Big format", format: 150},
	{name: "0 format", format: 0},
	{name: "Minus format", format: -1},
}

// TestGetWeatherWrongFormat checks that invalid format returns ErrWrongFormat.
func TestGetWeatherWrongFormat(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			geo := geo.GeoData{
				City: "London",
			}

			_, err := weather.GetWeather(geo, tc.format)
			if err != weather.ErrWrongFormat {
				t.Errorf("Expected error %v, got %v", weather.ErrWrongFormat, err)
			}
		})
	}

}
