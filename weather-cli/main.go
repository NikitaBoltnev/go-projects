// Package main is a CLI application that retrieves and displays weather information.
// It determines the user's location (by city or IP) and fetches weather from wttr.in.
package main

import (
	"flag"
	"fmt"
	"weather-cli/geo"
	"weather-cli/weather"
)

// main runs the weather CLI: parses flags, gets location, fetches and prints weather.
func main() {
	city := flag.String("city", "", "User's city")
	format := flag.Int("format", 1, "Output format for weather (1-4)")

	flag.Parse()
	geoData, err := geo.GetMyLocation(*city)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(geoData)

	weatherData, err := weather.GetWeather(*geoData, *format)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(weatherData)
}
