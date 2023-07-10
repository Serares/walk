package main

func filterEmptyStrings(list []string) []string {
	var filteredList []string
	for _, s := range list {
		if s != "" {
			filteredList = append(filteredList, s)
		}
	}
	return filteredList
}
