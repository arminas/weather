package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"
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
	Selected Location
	Forecast ForecastResponse
}

func (l Locations) Len() int           { return len(l) }
func (l Locations) Less(i, j int) bool { return l[i].Order < l[j].Order }
func (l Locations) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l Locations) LocationFromCity(city string) Location {
	for _, location := range l {
		if location.Name == city {
			return location
		}
	}

	return l[0]
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
	selectedCity := locations.LocationFromCity(selectedPlace)

	fmt.Printf("received request for %s\n", selectedPlace)

	sort.Sort(locations)

	var forecast = fetch_weather(selectedCity)

	templ_data := TemplateData{
		Options:  locations,
		Selected: selectedCity,
		Forecast: forecast,
	}

	t.Execute(w, templ_data)
}
