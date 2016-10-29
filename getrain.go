package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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

func CheckTime(t string) bool {
	currtime := time.Now()
	t1, _ := time.Parse(
		time.RFC3339,
		t)
	return currtime.Before(t1)
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
		log.Print(err)
		return ("Cloud'n svarer ikke :(")
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
		if CheckTime(data.Time) {
			fmt.Println(data.Time)
			rain := data.Precipitation.Intensity
			total += rain
		}
	}

	rain := int(total / float64(len(data.Points)))

	switch {
	case rain < 1:
		return ("Faktisk ikke")
	case rain >= 1 && rain <= 5:
		return ("Ja for faen")
	case rain > 5 && rain <= 20:
		return ("Omg ja, hold deg inne")
	default:
		return ("Usikker")
	}

}

func GetRaw() *Data {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	data := Data{}
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	return &data
}
