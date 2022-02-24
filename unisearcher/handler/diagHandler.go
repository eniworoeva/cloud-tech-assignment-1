package handler

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"time"
	"unisearcher/functions"
	"unisearcher/model"
)

// DiagHandler
//Diagnostic handler to showcases access to request content (headers, body, method, parameters, etc.)
func DiagHandler(w http.ResponseWriter, r *http.Request) {

	//Guard that prevents longer, unnecessary paths. ie if user enters /unisearcher/v1/diag/something/ it will return a http error
	if m, err := path.Match("/unisearcher/v1/diag/", r.URL.Path); !m && (err == nil) {
		http.Error(w, "Wrong usage of API. \nCorrect usage is /unisearcher/v1/diag/", http.StatusBadRequest)
		return
	}

	//Issues requests to the external apis
	uniDiag := functions.SendRequest("http://universities.hipolabs.com/")
	countryDiag := functions.SendRequest("https://restcountries.com/")

	//Initialize empty variables
	uniApiR, countryApiR := "", ""

	//Checks connection to Uni API
	if uniDiag == nil {
		uniApiR = strconv.Itoa(http.StatusBadGateway)
	} else {
		uniApiR = uniDiag.Status
	}

	//Checks connection to Country API
	if countryDiag == nil {
		countryApiR = strconv.Itoa(http.StatusBadGateway)
	} else {
		countryApiR = countryDiag.Status
	}

	uptime := time.Since(functions.GetUpTime()).Seconds()

	//Instantiates Diag
	d := model.Diag{
		UniversityAPI: uniApiR,
		CountryAPI:    countryApiR,
		Version:       VERSION,
		Uptime:        strconv.Itoa(int(uptime)) + "s",
	}

	// Write content type header
	w.Header().Add("content-type", "application/json")

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	//Encodes diag
	err := encoder.Encode(d)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}
}
