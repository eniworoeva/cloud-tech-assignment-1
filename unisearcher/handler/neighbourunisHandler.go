package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"unisearcher/functions"
	"unisearcher/model"
)

func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}
	//Retrieves necessary information from path
	cDir, rName := path.Split(r.URL.Path)

	limit := 10000

	// Space friendly query :)
	name := strings.Replace(rName, " ", "%20", -1)

	if validQuery := strings.Split(r.URL.RawQuery, "="); len(validQuery) != 2 && len(r.URL.RawQuery) != 0 {
		http.Error(w, "Bad usage of query..\nCorrect usage is /unisearcher/v1/neighbourunis/{country}/{name_or_partial_name}?limit={any_postive_number}", http.StatusBadRequest)
		return
	}

	if len(r.URL.RawQuery) != 0 {
		if _, err := strconv.Atoi(strings.Split(r.URL.RawQuery, "=")[1]); err != nil {
			http.Error(w, "Bad usage of query..\nCorrect usage is /unisearcher/v1/neighbourunis/{country}/{name_or_partial_name}?limit={any_postive_number}", http.StatusBadRequest)
			return
		}
		if validQuery := strings.Split(r.URL.RawQuery, "="); len(validQuery) != 2 {
			http.Error(w, "Bad usage of query..\nCorrect usage is /unisearcher/v1/neighbourunis/{country}/{name_or_partial_name}?limit={any_postive_number}", http.StatusBadRequest)
			return
		}
		if limit, _ = strconv.Atoi(strings.Split(r.URL.RawQuery, "=")[1]); limit < 1 {
			http.Error(w, "Limit should be greater than/equal to 1", http.StatusBadRequest)
			return
		}
	}

	if t := strings.Count(r.URL.Path, "/"); t != 5 {
		http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/neighbourunis/{country}/{name_or_partial_name}", http.StatusBadRequest)
		return
	}

	if len(name) == 0 {
		http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/neighbourunis/{country}/{name_or_partial_name}", http.StatusBadRequest)
		return
	}

	country := path.Base(cDir)
	url := fmt.Sprintf("https://restcountries.com/v3.1/name/%s", country)

	//Empty slice
	bordersCache := make([]model.BordersCache, 0)

	//Sends request to external API, returns JSON decoder
	borderRequest := functions.SendRequest(url)
	if borderRequest == nil {
		http.Error(w, "Error connecting to Country API", http.StatusBadGateway)
		return
	}

	decoder := json.NewDecoder(borderRequest.Body)
	//Decodes request, if successful -> continue.
	if err := decoder.Decode(&bordersCache); err != nil {
		http.Error(w, "No results found", http.StatusNotFound)
		return
	}

	fmt.Println(bordersCache)
	//Slice containing cca3 codes from country and bordering
	var b []string

	for i := 0; i < len(bordersCache); i++ {
		if !functions.Contains(b, bordersCache[i].CCA3) {
			b = bordersCache[i].Borders
			b = append(b, bordersCache[i].CCA3)
		}
	}

	fmt.Println(b)

	url = fmt.Sprintf("https://restcountries.com/v3.1/alpha?codes=%s", strings.Join(b[:], ","))

	countryRequest := functions.SendRequest(url)
	if countryRequest == nil {
		http.Error(w, "Error connecting to Country API", http.StatusBadGateway)
		return
	}

	decoder = json.NewDecoder(countryRequest.Body)

	//Creates empty slice
	countryCache := make([]model.CountryCache, 0)

	//Decodes request
	if err := decoder.Decode(&countryCache); err != nil {
		log.Fatal(err)
	}

	// URL to invoke
	url = fmt.Sprintf("http://universities.hipolabs.com/search?name=%s", name)

	uniRequest := functions.SendRequest(url)
	if uniRequest == nil {
		http.Error(w, "Error connecting to Uni API", http.StatusBadGateway)
		return
	}

	decoder = json.NewDecoder(uniRequest.Body)

	unis := make([]model.UniCache, 0)

	if err := decoder.Decode(&unis); err != nil {
		log.Fatal(err)
	}

	if len(unis) == 0 {
		http.Error(w, "No results found", http.StatusNotFound)
		return
	}

	//Creates an empy slice
	out := make([]model.UniInfoResponse, 0, limit)

	//Uses information from UniCache and CountryCache to create a new struct with the correct fields
	for _, obj := range unis {

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

		if len(out) == limit {
			break
		}
	}
	functions.EncodeUniInfo(w, out)
}
