package functions

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
