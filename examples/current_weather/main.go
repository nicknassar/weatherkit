package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nicknassar/weatherkit"
)

// print the current temp in new york
func main() {
	client, err := weatherkit.NewClient(
		os.Getenv("WEATHER_KIT_KID"),
		os.Getenv("WEATHER_KIT_ISS"),
		os.Getenv("WEATHER_KIT_SUB"),
		os.Getenv("WEATHER_KIT_PRIVATE_KEY"))

	if err != nil {
		fmt.Println("error", err.Error())
		return
	}

	ctx := context.Background()

	request := weatherkit.WeatherRequest{
		Latitude:  38.960,
		Longitude: -104.506,
		Language:  "en",
		DataSets: weatherkit.DataSets{
			weatherkit.DataSetCurrentWeather,
		},
	}

	weather, err := client.Weather(ctx, request)
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}

	fmt.Println(weather.CurrentWeather.Temperature)
}
