package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Location struct {
	Name      string
	Latitude  string
	Longitude string
	Order     int
}
type Locations []Location

type TemplateData struct {
	Options  Locations
	Selected string
	Forecast ForecastResponse
}

type ForecastResponse struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Hourly    struct {
		Time        []int64   `json:"time"`
		Temperature []float32 `json:"temperature_2m"`
		Humidity    []int     `json:"relative_humidity_2m"`
	} `json:"hourly"`
}

func (r ForecastResponse) Weather() []Weather {
	var list []Weather

	for index, timestamp := range r.Hourly.Time {
		parsedTime := time.Unix(timestamp, 0)

		list = append(list, Weather{
			DateTime:    parsedTime,
			Temperature: r.Hourly.Temperature[index],
			Humidity:    r.Hourly.Humidity[index],
		})
	}

	return list
}

type Weather struct {
	DateTime    time.Time
	Temperature float32
	Humidity    int
}

func (l Locations) Len() int           { return len(l) }
func (l Locations) Less(i, j int) bool { return l[i].Order < l[j].Order }
func (l Locations) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l Locations) LocationFromCity(city string) (Location, error) {
	for _, location := range l {
		if location.Name == city {
			return location, nil
		}
	}
	return Location{}, errors.New("city not found")
}

var locations = Locations{
	{Name: "Vilnius", Latitude: "54.68", Longitude: "25.27", Order: 0},
	{Name: "Kaunas", Latitude: "54.90", Longitude: "23.89", Order: 1},
	{Name: "Klaipėda", Latitude: "55.70", Longitude: "21.12", Order: 2},
	{Name: "Šiauliai", Latitude: "55.93", Longitude: "23.30", Order: 3},
	{Name: "Panevėžys", Latitude: "55.73", Longitude: "24.35", Order: 4},
}

func homepage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./ui/index.tmpl.html")

	selectedPlace := r.URL.Query().Get("place")
	if strings.TrimSpace(selectedPlace) == "" {
		selectedPlace = "Vilnius"
	}
	selectedCity, err := locations.LocationFromCity(selectedPlace)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("received request for %s\n", selectedPlace)

	sort.Sort(locations)

	var forecast = fetch_weather(selectedCity)

	templ_data := TemplateData{
		Options:  locations,
		Selected: selectedPlace,
		Forecast: forecast,
	}

	t.Execute(w, templ_data)
}

func fetch_weather(location Location) ForecastResponse {
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
	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("client: response body: %s\n", body)

	var parsed ForecastResponse
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		log.Fatalf("Unable to marshal JSON due to %s", err)
	}

	return parsed
}
