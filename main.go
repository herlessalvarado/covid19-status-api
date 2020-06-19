package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

// CovidWorldwide initial struct
type CovidWorldwide struct {
	Cases     string `json:"cases"`
	Deaths    string `json:"deaths"`
	Recovered string `json:"recovered"`
}

var covidWorldwide CovidWorldwide

// GetAll gets all worldwide numbers
func GetAll() {
	res, err := http.Get("https://www.worldometers.info/coronavirus/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".maincounter-number").Each(func(i int, s *goquery.Selection) {
		number := s.Find("span").Text()
		if i == 0 {
			covidWorldwide.Cases = number
		} else if i == 1 {
			covidWorldwide.Deaths = number
		} else {
			covidWorldwide.Recovered = number
		}
	})

	fmt.Println("Global cases updated")
}

func main() {
	GetAll()

	r := mux.NewRouter()

	r.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(covidWorldwide)
	}).Methods("GET")

	http.ListenAndServe(":8000", r)
}
