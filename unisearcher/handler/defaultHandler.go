package handler

import (
	"net/http"
)

// EmptyHandler /*
func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "No functionality on this level. Please use "+UNIINFO_PATH+", "+NEIGHBOURUNIS_PATH+" or "+DIAG_PATH+"."+
		"\nYou can also refer to the README at https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2022-workspace/ivann/assignment-1/-/blob/main/README.md for more information.", http.StatusOK)
}
