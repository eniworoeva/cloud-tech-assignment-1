package functions

import "unisearcher/model"

//Taken from https://gosamples.dev/slice-contains/
//A simple function to check if a slice contains a string element
//Used to prevent duplicates
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func StructContains(uniinfo []model.UniInfoResponse, str string) bool {
	for _, v := range uniinfo {
		if v.Name == str {
			return true
		}
	}
	return false
}
