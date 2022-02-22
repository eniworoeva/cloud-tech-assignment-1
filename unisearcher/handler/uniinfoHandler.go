package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"unisearcher/functions"
	"unisearcher/model"
)

func UniInfoHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not supported. Currently only GET supported.", http.StatusNotImplemented)
		return
	}
	//Initialize variables
	url, query := "", ""

	if m, err := path.Match("/unisearcher/v1/uniinfo/", r.URL.Path); m && (err == nil) {
		http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/uniinfo/{query}", http.StatusBadRequest)
		return
	}

	if t := strings.Count(r.URL.Path, "/"); t != 4 {
		http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/uniinfo/{query}", http.StatusBadRequest)
		return
	}
	//Space friendly search :)
	query = strings.Replace(path.Base(r.URL.Path), " ", "%20", -1)

	// URL to invoke
	url = fmt.Sprintf("http://universities.hipolabs.com/search?name=%s", query)

	//Issues new request
	uniRequest := functions.SendRequest(url)
	if uniRequest == nil {
		return
	}

	decoder := json.NewDecoder(uniRequest.Body)

	//Creates empty slice
	unis := make([]model.UniCache, 0)

	//Decodes request
	if err := decoder.Decode(&unis); err != nil {
		log.Fatal(err)
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
		return
	}

	decoder = json.NewDecoder(countryRequest.Body)

	//Creates empty slice
	countries := make([]model.CountryCache, 0)

	//Decodes request
	if err := decoder.Decode(&countries); err != nil {
		log.Fatal(err)
	}

	//Creates an empty slice
	out := make([]model.UniInfoResponse, 0, len(unis))

	//Uses information from UniCache and CountryCache to create a new struct with the correct fields
	for _, obj := range unis {
		for _, c := range countries {
			if c.CCA2 == obj.AlphaTwoCode && (!functions.StructContains(out, obj.Name)) {
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
	//Return response. See defaultHandler for method
	functions.EncodeUniInfo(w, out)
}
