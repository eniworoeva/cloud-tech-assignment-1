package functions

import "unisearcher/model"

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

// StructContains
/* Variation of method above
A simple function to check if a slice contains an element of UniInfoResponse
Used to prevent duplicates */
func StructContains(uniinfo []model.UniInfoResponse, str string) bool {
	for _, v := range uniinfo {
		if v.Name == str {
			return true
		}
	}
	return false
}
