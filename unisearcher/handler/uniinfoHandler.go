package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"unisearcher/functions"
)

// UniInfoHandler /*
func UniInfoHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		handleGetRequest(w, r)
	default:
		http.Error(w, "Method not supported. Currently only GET are supported.", http.StatusNotImplemented)
		return
	}
}

/*
Dedicated handler for GET requests
*/
func handleGetRequest(w http.ResponseWriter, r *http.Request) {

	// Space friendly query :)
	query := strings.Replace(path.Base(r.URL.Path), " ", "%20", -1)

	// URL to invoke
	url := fmt.Sprintf("http://universities.hipolabs.com/search?name=%s", query)

	// Instantiate the client
	client := &http.Client{}

	// Create new request
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Errorf("Error in creating request:", err.Error())
	}

	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		fmt.Errorf("Error in response", err.Error())
	}

	decoder := json.NewDecoder(res.Body)

	UniCache := make([]UniCache, 0)
	if err := decoder.Decode(&UniCache); err != nil {
		log.Fatal(err)
	}
	// Creates empty slice
	cca2 := make([]string, 0)
	for _, uni := range UniCache {
		if !functions.Contains(cca2, uni.AlphaTwoCode) {
			cca2 = append(cca2, uni.AlphaTwoCode)
		}
	}
	_ = fmt.Sprintf("https://restcountries.com/v3.1/alpha?codes=%s", strings.Join(cca2[:], ","))

	out := make([]Response, 0, len(UniCache))
	for _, obj := range UniCache {
		out = append(out, Response{Name: obj.Name, Country: obj.Country, IsoCode: obj.AlphaTwoCode, WebPages: obj.WebPages})
	}

	sendResponse(w, out)
}

func sendResponse(w http.ResponseWriter, r []Response) {
	// Write content type header
	w.Header().Add("content-type", "application/json")

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	err := encoder.Encode(r)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}
}
