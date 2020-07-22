package main

import (
	"fmt"

	"github.com/shawntoffel/darksky"
)

func main() {
	fmt.Println(GetDarkRain())
}

// GetDarkRain lol
func GetDarkRain() string {
	client := darksky.New("5acf2ebae35368eddd649a78e7ef27dc")

	request := darksky.ForecastRequest{}
	request.Latitude = 60.3919
	request.Longitude = 5.3228
	request.Options = darksky.ForecastRequestOptions{Exclude: "hourly,minutely"}

	response, err := client.Forecast(request)

	if err != nil {
		fmt.Println(err.Error())
		return "Clouden svarer ikke"
	}

	// fmt.Println(response.Currently.PrecipIntensity)
	// fmt.Println(response.Currently.PrecipProbability)

	switch {
	case response.Currently.PrecipProbability < 0.2:
		return ("Faktisk ikke")
	case response.Currently.PrecipProbability >= 0.5 && response.Currently.PrecipProbability <= 0.8:
		return ("Ja for faen")
	case response.Currently.PrecipProbability > 1.0 && response.Currently.PrecipProbability <= 20.0:
		return ("Omg ja, hold deg inne")
	default:
		return ("Usikker")
	}
}
