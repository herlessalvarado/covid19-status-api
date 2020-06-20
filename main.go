package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
)

// CovidWorldwide initial struct
type CovidWorldwide struct {
	Cases     int `json:"cases"`
	Deaths    int `json:"deaths"`
	Recovered int `json:"recovered"`
}

var covidWorldwide CovidWorldwide

// CountryCovid initial struct
type CountryCovid struct {
	Country        string `json:"country"`
	TotalCases     int    `json:"totalCases"`
	NewCases       int    `json:"newCases"`
	TotalDeaths    int    `json:"totalDeaths"`
	NewDeaths      int    `json:"newDeaths"`
	TotalRecovered int    `json:"totalRecovered"`
	ActiveCases    int    `json:"activeCases"`
	CriticalCases  int    `json:"criticalCases"`
	TotalTests     int    `json:"totalTests"`
	Population     int    `json:"population"`
}

var countryCovid []CountryCovid

var replacer = strings.NewReplacer(" ", "", ",", "", "+", "")

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
			number = replacer.Replace(number)
			i1, _ := strconv.Atoi(number)
			covidWorldwide.Cases = i1
		} else if i == 1 {
			number = replacer.Replace(number)
			i1, _ := strconv.Atoi(number)
			covidWorldwide.Deaths = i1
		} else {
			number = replacer.Replace(number)
			i1, _ := strconv.Atoi(number)
			covidWorldwide.Recovered = i1
		}
	})

	fmt.Println("Global numbers updated")
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
			countryName := s.Text()
			var country CountryCovid
			country.Country = countryName
			countryCovid = append(countryCovid, country)
		}
		if i%countriesTableLenght == 2 {
			totalCases := s.Text()
			totalCases = replacer.Replace(totalCases)
			i1, _ := strconv.Atoi(totalCases)
			countryCovid[len(countryCovid)-1].TotalCases = i1
		}
		if i%countriesTableLenght == 3 {
			newCases := s.Text()
			newCases = replacer.Replace(newCases)
			i1, _ := strconv.Atoi(newCases)
			countryCovid[len(countryCovid)-1].NewCases = i1
		}
		if i%countriesTableLenght == 4 {
			totalDeaths := s.Text()
			totalDeaths = replacer.Replace(totalDeaths)
			i1, _ := strconv.Atoi(totalDeaths)
			countryCovid[len(countryCovid)-1].TotalDeaths = i1
		}
		if i%countriesTableLenght == 5 {
			newDeaths := s.Text()
			newDeaths = replacer.Replace(newDeaths)
			i1, _ := strconv.Atoi(newDeaths)
			countryCovid[len(countryCovid)-1].NewDeaths = i1
		}
		if i%countriesTableLenght == 6 {
			totalRecovered := s.Text()
			totalRecovered = replacer.Replace(totalRecovered)
			i1, _ := strconv.Atoi(totalRecovered)
			countryCovid[len(countryCovid)-1].TotalRecovered = i1
		}
		if i%countriesTableLenght == 8 {
			activeCases := s.Text()
			activeCases = replacer.Replace(activeCases)
			i1, _ := strconv.Atoi(activeCases)
			countryCovid[len(countryCovid)-1].ActiveCases = i1
		}
		if i%countriesTableLenght == 9 {
			criticalCases := s.Text()
			criticalCases = replacer.Replace(criticalCases)
			i1, _ := strconv.Atoi(criticalCases)
			countryCovid[len(countryCovid)-1].CriticalCases = i1
		}
		if i%countriesTableLenght == 12 {
			totalTests := s.Text()
			totalTests = replacer.Replace(totalTests)
			i1, _ := strconv.Atoi(totalTests)
			countryCovid[len(countryCovid)-1].TotalTests = i1
		}
		if i%countriesTableLenght == 14 {
			population := s.Text()
			population = replacer.Replace(population)
			i1, _ := strconv.Atoi(population)
			countryCovid[len(countryCovid)-1].Population = i1
		}
	})

	countryCovid = countryCovid[1:]

	if len(countryCovid) > 0 {
		countryCovid = countryCovid[:len(countryCovid)-1]
	}

	sort.SliceStable(countryCovid, func(i, j int) bool {
		return countryCovid[i].TotalCases > countryCovid[j].TotalCases
	})

	fmt.Println("Country cases updated")
}

func main() {
	GetAll()
	GetCountries()

	c := cron.New()
	c.AddFunc("@every 5m", GetAll)
	c.AddFunc("@every 5m", GetCountries)
	c.Start()

	r := mux.NewRouter()

	r.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(covidWorldwide)
	}).Methods("GET")

	r.HandleFunc("/countries", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(countryCovid)
	}).Methods("GET")

	r.HandleFunc("/countries/{country}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		country := vars["country"]
		var requestedCountry CountryCovid
		for _, v := range countryCovid {
			if v.Country == country {
				requestedCountry = v
			}
		}
		json.NewEncoder(w).Encode(requestedCountry)
	}).Methods("GET")

	http.ListenAndServe(":8000", r)
}
