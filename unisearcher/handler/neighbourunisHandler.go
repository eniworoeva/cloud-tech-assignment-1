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

// NeighbourUnisHandler /*
func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {
	//Error guard that prohibts requests that are not of type GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}
	//Retrieves necessary information from path
	cDir, rName := path.Split(r.URL.Path)

	//Initialize limit beforehand if no limit is used
	limit := 10000

	// Space friendly search :)
	name := strings.Replace(rName, " ", "%20", -1)

	//Error guarding against wrongful usage of query
	if len(r.URL.RawQuery) != 0 {
		if _, err := strconv.Atoi(strings.Split(r.URL.RawQuery, "=")[1]); err != nil {
			http.Error(w, "Limit should be a number..\nCorrect usage is /unisearcher/v1/neighbourunis/{:country_name}/{:partial_or_complete_university_name}{?limit={:number}}", http.StatusBadRequest)
			return
		}
		if validQuery := strings.Split(r.URL.RawQuery, "="); len(validQuery) != 2 {
			http.Error(w, "Wrong usage of query..\nCorrect usage is /unisearcher/v1/neighbourunis/{:country_name}/{:partial_or_complete_university_name}{?limit={:number}}", http.StatusBadRequest)
			return
		}
		if limit, _ = strconv.Atoi(strings.Split(r.URL.RawQuery, "=")[1]); limit < 1 {
			http.Error(w, "Limit should be greater than/equal to 1", http.StatusBadRequest)
			return
		}
	}

	//Error guard that prohibits blank search-parameter
	if m, err := path.Match("/unisearcher/v1/neighbourunis/", r.URL.Path); m && (err == nil) {
		http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/neighbourunis/{:country_name}/{:partial_or_complete_university_name}{?limit={:number}}", http.StatusBadRequest)
		return
	}

	//Error guard that prohibits use of wrong path. ie /unisearcher/v1/neighbourunis/country/test/test"
	if t := strings.Count(r.URL.Path, "/"); t != 5 {
		http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/neighbourunis/{:country_name}/{:partial_or_complete_university_name}{?limit={:number}}", http.StatusBadRequest)
		return
	}

	//Gets country
	country := path.Base(cDir)

	//Link to API
	url := fmt.Sprintf("https://restcountries.com/v3.1/name/%s", country)

	//Empty slice
	bordersCache := make([]model.BordersCache, 0)

	//Issues new request
	borderRequest := functions.SendRequest(url)
	if borderRequest == nil {
		http.Error(w, "Error connecting to Country API", http.StatusBadGateway)
		return
	}

	//Creates decoder to decode response from GET request
	decoder := json.NewDecoder(borderRequest.Body)

	//Populates slice if decoder is successful
	if err := decoder.Decode(&bordersCache); err != nil {
		http.Error(w, "No results found. There are either no countries with that name, or you've used an unsupported character", http.StatusNotFound)
		return
	}

	//Slice containing cca3 codes from country and bordering
	var b []string

	//Loops over all bordering country codes in borderscache and adds them to a slice, ignores duplicates.
	//Some country name searches returns multiple countries. ie using "pri" as a country name parameter,
	//would result in a response with countries containing principality in their name.
	//This loop makes sure to add all countries and their bordering countries to the slice.
	for i := 0; i < len(bordersCache); i++ {
		tempArr := bordersCache[i].Borders
		if len(tempArr) == 0 && !functions.Contains(b, bordersCache[i].CCA3) {
			b = append(b, bordersCache[i].CCA3)
		}

		for j := 0; j < len(tempArr); j++ {
			if cca3 := bordersCache[i].CCA3; !functions.Contains(b, cca3) {
				b = append(b, cca3)
			}
			b = append(b, tempArr[j])

		}
	}

	//New url to invoke, using cca3 list as parameter
	url = fmt.Sprintf("https://restcountries.com/v3.1/alpha?codes=%s", strings.Join(b[:], ","))

	//Issues new request
	countryRequest := functions.SendRequest(url)
	if countryRequest == nil {
		http.Error(w, "Error connecting to Country API", http.StatusBadGateway)
		return
	}
	//Creates decoder to decode response from GET request
	decoder = json.NewDecoder(countryRequest.Body)

	//Creates empty slice. Will contain information about all the countries in the cca3 slice
	countryCache := make([]model.CountryCache, 0)

	//Populates slice if decoder is successful
	if err := decoder.Decode(&countryCache); err != nil {
		log.Fatal(err)
	}

	// Last url to invoke
	url = fmt.Sprintf("http://universities.hipolabs.com/search?name_contains=%s", name)

	//Issues new request
	uniRequest := functions.SendRequest(url)
	if uniRequest == nil {
		http.Error(w, "Error connecting to Uni API", http.StatusBadGateway)
		return
	}

	//Creates decoder to decode response from GET request
	decoder = json.NewDecoder(uniRequest.Body)

	//Creates empty slice
	unis := make([]model.UniCache, 0)

	//Populates slice if decoder is successful
	if err := decoder.Decode(&unis); err != nil {
		log.Fatal(err)
	}

	//Guards against empty slice of unis, prevents unnecessary use of EncodeUniInfo to encode empty slice
	if len(unis) == 0 {
		http.Error(w, "No results found", http.StatusNotFound)
		return
	}

	//Creates an empy slice
	out := make([]model.UniInfoResponse, 0, limit)

	//Uses information from UniCache and CountryCache to create a new struct with the correct fields
	for _, obj := range unis {
		for _, c := range countryCache {
			if c.CCA2 == obj.AlphaTwoCode {
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
		//If there is a limit, stop loop when limit is reached
		if len(out) == limit {
			break
		}
	}

	//Guards against empty slice of outgoing slice, prevents unnecessary use of EncodeUniInfo to encode the empty slice
	if len(out) == 0 {
		http.Error(w, "No results found", http.StatusNotFound)
		return
	}

	//Encodes the outgoing slice
	functions.EncodeUniInfo(w, out)
}
