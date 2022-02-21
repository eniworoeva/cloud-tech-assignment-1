package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unisearcher/model"
)

/*
Empty handler
*/
func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "No functionality on this level. Please use "+UNIINFO_PATH+", "+NEIGHBOURUNIS_PATH+" or "+DIAG_PATH, http.StatusOK)
}

//Method used to write response to API
func sendResponse(w http.ResponseWriter, r []model.UniInfoResponse) {
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

//Method used to send GET requests to external APIs. Returns JSON decoder
func sendRequest(url string) *json.Decoder {
	// Create new request
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		_ = fmt.Errorf("error in creating request", err.Error())
	}

	// Instantiate the client
	client := &http.Client{}

	// Setting content type -> effect depends on the service provider
	r.Header.Add("content-type", "application/json")

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		_ = fmt.Errorf("Error in response", err.Error())
	}

	decoder := json.NewDecoder(res.Body)

	return decoder
}
