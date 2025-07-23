// Package geo_test contains tests for the geo package.
package geo_test

import (
	"testing"
	"weather-cli/geo"
)

// TestGetMyLocation verifies that a valid city returns correct GeoData.
func TestGetMyLocation(t *testing.T) {
	city := "London"

	expected := geo.GeoData{
		City: "London",
	}

	got, err := geo.GetMyLocation(city)

	if err != nil {
		t.Error(err)
	}

	if expected.City != got.City {
		t.Errorf("Expected %v, got %v", expected, got)
	}

}

// TestGetMyLocationNoCity checks that an invalid city returns ErrNoCity.
func TestGetMyLocationNoCity(t *testing.T) {
	city := "Londonssss"

	_, err := geo.GetMyLocation(city)
	if err != geo.ErrNoCity {
		t.Errorf("Expected error %v, got %v", geo.ErrNoCity, err)
	}
}
