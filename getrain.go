package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const url = "https://www.yr.no/api/v0/locations/id/1-92416/forecast/now"

type Data struct {
	_links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Created string `json:"created"`
	Points  []struct {
		Precipitation struct {
			Intensity float64 `json:"intensity"`
		} `json:"precipitation"`
		Time string `json:"time"`
	} `json:"points"`
	Update string `json:"update"`
}

func GetRain() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	data := Data{}
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	var total float64 = 0.0

	for _, data := range data.Points {
		rain := data.Precipitation.Intensity
		total += rain
	}

	rain := int(total / float64(len(data.Points)))

	switch {
	case rain < 1:
		return ("Regner ikke")
	case rain >= 1 && rain <= 5:
		return ("Regner")
	case rain > 5 && rain <= 20:
		return ("Heavy rain")
	default:
		return ("No hit")
	}

}
