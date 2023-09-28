package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"osaigie/unisearcher/functions"
	"osaigie/unisearcher/model"
	"path"
	"strings"
)



func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	}




// UniInfoHandler /*
func UniInfoHandler(w http.ResponseWriter, r *http.Request) {
//handle cors
	enableCors(&w)



	//Error guard that prohibts requests that are not of type GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}

	//Initialize variables
	url, search := "", ""

	//Error guard that prohibits blank search-parameter
	if m, err := path.Match("/unisearcher/v1/uniinfo/", r.URL.Path); m && (err == nil) {
		http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/uniinfo/{:partial_or_complete_university_name}/", http.StatusBadRequest)
		return
	}

	//Error guard that prohibits use of wrong path. ie /unisearcher/v1/uniinfo/example/test"
	if t := strings.Count(r.URL.Path, "/"); t != 4 {
		http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/uniinfo/{:partial_or_complete_university_name}/", http.StatusNotFound)
		return
	}

	//Error guard that prohibits use of optional queries as this endpoint is not supposed to have any
	if len(r.URL.RawQuery) != 0 {
		http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/uniinfo/{:partial_or_complete_university_name}/", http.StatusBadRequest)
		return
	}

	//Space friendly search :)
	search = strings.Replace(path.Base(r.URL.Path), " ", "%20", -1)

	// URL to invoke
	url = fmt.Sprintf("http://universities.hipolabs.com/search?name_contains=%s", search)

	//Issues new request
	uniRequest := functions.SendRequest(url)
	if uniRequest == nil {
		http.Error(w, "Error connecting to Uni API", http.StatusBadGateway)
		return
	}

	//Creates decoder to decode response from GET request
	decoder := json.NewDecoder(uniRequest.Body)

	//Creates empty slice
	unis := make([]model.UniCache, 0)

	//Populates slice if decoder is successful
	if err := decoder.Decode(&unis); err != nil {
		log.Fatal(err)
	}

	//Prevents unnecessary usage of REST countries api by guarding against empty responses
	if len(unis) == 0 {
		http.Error(w, "No results found", http.StatusNotFound)
		return
	}

	//Creates empty Slice
	cca2 := make([]string, 0)

	// Adds isocodes into the slice, ignores duplicates
	for _, uni := range unis {
		if !functions.Contains(cca2, uni.AlphaTwoCode) {
			cca2 = append(cca2, uni.AlphaTwoCode)
		}
	}

	//New url
	url = fmt.Sprintf("https://restcountries.com/v3.1/alpha?codes=%s", strings.Join(cca2[:], ","))

	//Issues new request
	countryRequest := functions.SendRequest(url)
	if countryRequest == nil {
		http.Error(w, "Error connecting to Country API", http.StatusBadGateway)
		return
	}

	//Creates decoder to decode response from GET request
	decoder = json.NewDecoder(countryRequest.Body)

	//Creates empty slice
	countries := make([]model.CountryCache, 0)

	//Populates slice if decoder is successful
	if err := decoder.Decode(&countries); err != nil {
		log.Fatal(err)
	}

	//Creates an empty slice
	out := make([]model.UniInfoResponse, 0, len(unis))

	//Uses information from UniCache and CountryCache to create a new struct with the combined fields of the other structs
	for _, obj := range unis {
		for _, c := range countries {
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
	}

	//Encodes the outgoing slice
	functions.EncodeUniInfo(w, out)
}
