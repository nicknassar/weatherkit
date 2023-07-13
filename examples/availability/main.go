package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nicknassar/weatherkit"
)

// print data set availability in new york
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

	availability, err := client.Availability(ctx, weatherkit.AvailabilityRequest{
		Latitude:  38.960,
		Longitude: -104.506,
	})
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}

	fmt.Println(availability)
}
