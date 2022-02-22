package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"unisearcher/functions"
	"unisearcher/model"
)

/*
 Diagnostic handler to showcases access to request content (headers, body, method, parameters, etc.)
*/
func DiagHandler(w http.ResponseWriter, r *http.Request) {
	uniDiag := functions.SendRequest("http://universities.hipolabs.com/")
	countryDiag := functions.SendRequest("https://restcountries.com/")

	uniApiR, countryApiR := "", ""

	if uniDiag == nil {
		uniApiR = strconv.Itoa(http.StatusNotFound)
	} else {
		uniApiR = uniDiag.Status
	}

	if countryDiag == nil {
		countryApiR = strconv.Itoa(http.StatusNotFound)
	} else {
		countryApiR = countryDiag.Status
	}

	d := model.Diag{
		UniversityAPI: uniApiR,
		CountryAPI:    countryApiR,
		Version:       VERSION,
		Uptime:        fmt.Sprint(time.Since(functions.GetUpTime()).Round(time.Second)),
	}

	// Write content type header
	w.Header().Add("content-type", "application/json")

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	//Encodes response
	err := encoder.Encode(d)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}

}
