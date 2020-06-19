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

// CountryCovid initial struct
type CountryCovid struct {
	Name           string `json:"name"`
	TotalCases     string `json:"totalCases"`
	NewCases       string `json:"newCases"`
	TotalDeaths    string `json:"totalDeaths"`
	NewDeaths      string `json:"newDeaths"`
	TotalRecovered string `json:"totalRecovered"`
	ActiveCases    string `json:"activeCases"`
	CriticalCases  string `json:"criticalCases"`
	TotalTests     string `json:"totalTests"`
	Population     string `json:"population"`
}

var countryCovid []CountryCovid

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

// GetCountries gets numbers for each country
func GetCountries() {
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

	var countriesTable = doc.Find("table#main_table_countries_today")
	var countriesTableLenght = doc.Find("table#main_table_countries_today th").Length()
	countriesTable.Find("tbody").Find("tr:not(.row_continent)").Find("td").Each(func(i int, s *goquery.Selection) {
		if i%countriesTableLenght == 1 {
			name := s.Text()
			var country CountryCovid
			country.Name = name
			countryCovid = append(countryCovid, country)
		}
		if i%countriesTableLenght == 2 {
			totalCases := s.Text()
			countryCovid[len(countryCovid)-1].TotalCases = totalCases
		}
		if i%countriesTableLenght == 3 {
			newCases := s.Text()
			countryCovid[len(countryCovid)-1].NewCases = newCases
		}
		if i%countriesTableLenght == 4 {
			totalDeaths := s.Text()
			countryCovid[len(countryCovid)-1].TotalDeaths = totalDeaths
		}
		if i%countriesTableLenght == 5 {
			newDeaths := s.Text()
			countryCovid[len(countryCovid)-1].NewDeaths = newDeaths
		}
		if i%countriesTableLenght == 6 {
			totalRecovered := s.Text()
			countryCovid[len(countryCovid)-1].TotalRecovered = totalRecovered
		}
		if i%countriesTableLenght == 8 {
			activeCases := s.Text()
			countryCovid[len(countryCovid)-1].ActiveCases = activeCases
		}
		if i%countriesTableLenght == 9 {
			criticalCases := s.Text()
			countryCovid[len(countryCovid)-1].CriticalCases = criticalCases
		}
		if i%countriesTableLenght == 12 {
			totalTests := s.Text()
			countryCovid[len(countryCovid)-1].TotalTests = totalTests
		}
		if i%countriesTableLenght == 14 {
			population := s.Text()
			countryCovid[len(countryCovid)-1].Population = population
		}
	})

	countryCovid = countryCovid[1:]

	if len(countryCovid) > 0 {
		countryCovid = countryCovid[:len(countryCovid)-1]
	}

	fmt.Println("Country cases updated")
}

func main() {
	GetAll()
	GetCountries()

	r := mux.NewRouter()

	r.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(covidWorldwide)
	}).Methods("GET")

	r.HandleFunc("/countries", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(countryCovid)
	}).Methods("GET")

	http.ListenAndServe(":8000", r)
}
