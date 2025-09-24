package main

import (
	"strings"
)

func NormalizeCSV(s string) string {
	parts := strings.Fields(s)

	var res []string
	for _, part := range parts {
		good := strings.Trim(part, ",")
		if good != "" {
			res = append(res, good)
		}
	}

	return strings.Join(res, ",")
}
