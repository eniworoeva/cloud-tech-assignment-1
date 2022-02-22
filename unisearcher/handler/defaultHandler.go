package handler

import (
	"net/http"
)

/*
Empty handler
*/
func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "No functionality on this level. Please use "+UNIINFO_PATH+", "+NEIGHBOURUNIS_PATH+" or "+DIAG_PATH, http.StatusOK)
}
