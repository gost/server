package odata

import (
	"strings"
)

var keywords = []string{"$filter", "$select", "$expand", "$orderby", "$top", "$skip", "$count"}

func partHasKeyword(part string) bool {
	part1 := strings.ToLower(strings.Split(part, "=")[0])

	for _, kw := range keywords {
		if part1 == kw {
			return true
		}
	}
	return false
}

// IsValidOdataQuery checks a querystring for containing correct keywords
func IsValidOdataQuery(query string) bool {
	res := true
	parts := strings.Split(query, "&")
	for _, element := range parts {
		if !partHasKeyword(element) {
			return false
		}
	}
	return res
}
