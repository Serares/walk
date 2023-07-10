package main

import (
	"fmt"
	"time"
)

const LayoutDate string = "2006-01-02"
const LayoutTime string = "15:04:05" // TODO support time

func filterEmptyStrings(list []string) []string {
	var filteredList []string
	for _, s := range list {
		if s != "" {
			filteredList = append(filteredList, s)
		}
	}
	return filteredList
}

// this will help validate the date provided by user
func validateDate(date string) string {
	_, err := time.Parse(LayoutDate, date)
	if err != nil {
		fmt.Println("Provided date is crap, replacing the provided date with an empty string")
		return ""
	}
	return date
}
