package main

import (
	"html/template"
	"net/http"
	"sort"
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
}

func (l Locations) Len() int           { return len(l) }
func (l Locations) Less(i, j int) bool { return l[i].Order < l[j].Order }
func (l Locations) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func homepage(w http.ResponseWriter, r *http.Request) {
	var locations = Locations{
		{Name: "Vilnius", Latitude: "54.68", Longitude: "25.27", Order: 0},
		{Name: "Kaunas", Latitude: "54.90", Longitude: "23.89", Order: 1},
		{Name: "Klaipėda", Latitude: "55.70", Longitude: "21.12", Order: 2},
		{Name: "Šiauliai", Latitude: "55.93", Longitude: "23.30", Order: 3},
		{Name: "Panevėžys", Latitude: "55.73", Longitude: "24.35", Order: 4},
	}

	t, _ := template.ParseFiles("./ui/index.tmpl.html")

	selected_place := r.URL.Query().Get("place")

	sort.Sort(locations)

	templ_data := TemplateData{
		Options:  locations,
		Selected: selected_place,
	}

	t.Execute(w, templ_data)
}