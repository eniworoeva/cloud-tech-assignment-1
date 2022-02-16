package handler

import "net/http"

/*
Empty handler
*/
func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "No functionality on root level", http.StatusOK)
}
