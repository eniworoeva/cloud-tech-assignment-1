package handler

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"unisearcher/functions"
	"unisearcher/model"
)

//UniInfoHandler
func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		//Retrieves necessary information from path
		fullPath, search := path.Split(r.URL.Path)
		country := path.Base(fullPath)
		query := r.URL.RawQuery

		url := fmt.Sprintf("https://restcountries.com/v3.1/name/%s", country)

		//Empty slice
		bordersCache := make([]model.BordersCache, 0)

		//Sends request to external API, returns JSON decoder
		borderRequest := sendRequest(url)

		//Decodes request, if successful -> continue.
		if err := borderRequest.Decode(&bordersCache); err != nil {
			http.Error(w, "No results", http.StatusNotFound)
		} else {
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
			url = fmt.Sprintf("http://universities.hipolabs.com/search?name=%s", search)

			uniRequest := sendRequest(url)

			unis := make([]model.UniCache, 0)
			if err := uniRequest.Decode(&unis); err != nil {
				log.Fatal(err)
			}

			capacity := len(unis)

			if query != "" {
				q := strings.Split(query, "=")

				limit, _ := strconv.Atoi(q[1])

				if limit > 0 {
					capacity = limit
				}
			}

			//Creates an empy slice
			out := make([]model.UniInfoResponse, 0, capacity)

			//Uses information from UniCache and CountryCache to create a new struct with the correct fields
			for _, obj := range unis {

				if len(out) == capacity {
					break
				}

				for _, c := range countryCache {
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
}
