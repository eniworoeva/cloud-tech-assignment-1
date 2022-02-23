package functions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"unisearcher/model"
)

// Initialize global variable
var timeStart = time.Now()

// GetUpTime Returns timeStart
func GetUpTime() time.Time {
	return timeStart
}

// Contains
/*Taken from https:gosamples.dev/slice-contains/
A simple function to check if a slice contains a string element
Used to prevent duplicates */
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// EncodeUniInfo Method used to write response to API, only used for neighbourunis and uniinfo api.
func EncodeUniInfo(w http.ResponseWriter, r []model.UniInfoResponse) {
	// Write content type header
	w.Header().Add("content-type", "application/json")

	// Instantiate encoder
	encoder := json.NewEncoder(w)

	//Encodes response
	err := encoder.Encode(r)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}
}

// SendRequest Method used to send GET requests to external APIs
func SendRequest(url string) *http.Response {
	// Create new request
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		_ = fmt.Errorf("error in creating request", err.Error())
	}

	// Instantiate the client
	client := &http.Client{}

	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")

	// Issues request
	res, err := client.Do(r)
	if err != nil || res == nil {
		_ = fmt.Errorf("Error in response", err.Error())
	}

	return res
}
