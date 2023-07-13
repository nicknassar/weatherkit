package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nicknassar/weatherkit"
)

// print event text for an alert id
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

	request := weatherkit.WeatherAlertRequest{
		ID:       "alert id",
		Language: "en",
	}

	response, err := client.Alert(ctx, request)
	if err != nil {
		fmt.Println("error", err.Error())
		return
	}

	fmt.Println(response.EventText[0].Text)
}
