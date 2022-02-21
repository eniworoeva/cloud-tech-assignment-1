package handler

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"unisearcher/functions"
	"unisearcher/model"
)

func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		//Retrieves necessary information from path
		fPath := "/unisearcher/v1/neighbourunis/"
		cDir, q := path.Split(r.URL.Path)
		country := ""

		if t := strings.Count(r.URL.Path, "/"); t != 5 {
			http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/neighbourunis/{country}/{query}", http.StatusBadRequest)
			return
		}

		if m, err := path.Match(fPath, r.URL.Path); m && (err == nil) {
			http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/neighbourunis/{country}/{query}", http.StatusBadRequest)
			return
		} else if m2, err2 := path.Match(fPath+q, r.URL.Path); m2 && (err2 == nil) {
			http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/neighbourunis/{country}/{query}", http.StatusBadRequest)
			return
		}

		if len(q) == 0 {
			http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/neighbourunis/{country}/{query}", http.StatusBadRequest)
			return
		}

		country = path.Base(cDir)
		url := fmt.Sprintf("https://restcountries.com/v3.1/name/%s", country)

		//Empty slice
		bordersCache := make([]model.BordersCache, 0)

		//Sends request to external API, returns JSON decoder
		borderRequest := sendRequest(url)

		//Decodes request, if successful -> continue.
		if err := borderRequest.Decode(&bordersCache); err != nil {
			log.Fatal(err)
		}

		//Slice containing cca3 codes from country and bordering countries
		b := bordersCache[0].Borders
		b = append(b, bordersCache[0].CCA3)

		url = fmt.Sprintf("https://restcountries.com/v3.1/alpha?codes=%s", strings.Join(b[:], ","))

		countryRequest := sendRequest(url)

		//Creates empty slice
		countryCache := make([]model.CountryCache, 0)

		//Decodes request
		if err := countryRequest.Decode(&countryCache); err != nil {
			log.Fatal(err)
		}

		// URL to invoke
		url = fmt.Sprintf("http://universities.hipolabs.com/search?name=%s", q)

		uniRequest := sendRequest(url)

		unis := make([]model.UniCache, 0)

		if err := uniRequest.Decode(&unis); err != nil {
			log.Fatal(err)
		}

		if len(unis) == 0 {
			http.Error(w, "No results found", http.StatusNotFound)
			return
		}

		//Creates an empy slice
		out := make([]model.UniInfoResponse, 0, len(unis))

		//Uses information from UniCache and CountryCache to create a new struct with the correct fields
		for _, c := range countryCache {
			for _, obj := range unis {
				if c.CCA2 == obj.AlphaTwoCode && !functions.StructContains(out, obj.Name) {
					out = append(out, model.UniInfoResponse{
						Name:      obj.Name,
						Country:   obj.Country,
						IsoCode:   obj.AlphaTwoCode,
						WebPages:  obj.WebPages,
						Languages: c.Languages,
						Map:       c.Maps["openStreetMaps"],
					})
				}
			}
		}
		sendResponse(w, out)
	}
}
