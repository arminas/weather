package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type ForecastResponse struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Hourly    struct {
		Time        []int64   `json:"time"`
		Temperature []float32 `json:"temperature_2m"`
		Humidity    []int     `json:"relative_humidity_2m"`
	} `json:"hourly"`
}

func (r ForecastResponse) Weather() Forecast {
	var forecast Forecast
	now := time.Now()

	for index, timestamp := range r.Hourly.Time {
		parsedTime := time.Unix(timestamp, 0)

		if parsedTime.Before(now) {
			continue
		}

		forecast.hourly = append(forecast.hourly, Weather{
			DateTime:    parsedTime,
			Temperature: r.Hourly.Temperature[index],
			Humidity:    r.Hourly.Humidity[index],
		})
	}

	return forecast
}

type Weather struct {
	DateTime    time.Time
	Temperature float32
	Humidity    int
}

type Forecast struct {
	hourly []Weather
}

func fetch_weather(location Location) Forecast {
	timezone := "Europe/Vilnius"
	requestURL := fmt.Sprintf(
		`https://api.open-meteo.com/v1/forecast?`+
			`latitude=%v`+
			`&longitude=%v`+
			`&hourly=temperature_2m,relative_humidity_2m`+
			`&timezone=%v`+
			`&timeformat=unixtime`+
			`&forecast_days=7`,
		location.Latitude, location.Longitude, timezone,
	)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "meno_lt_weather_scraper")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("client: got response!\n")
	log.Printf("client: status code: %d\n", res.StatusCode)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("client: response body: %s\n", body)

	var response ForecastResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("Unable to marshal JSON due to %s", err)
	}

	return response.Weather()
}
