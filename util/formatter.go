package util

import (
	"fmt"
	"strings"
)

// FormatDate formats date from YYYY-MM to Month YYYY
// example: 2020-10 -> October 2020
func FormatDate(date string) string {
	subStrings := strings.Split(date, "-")
	year := subStrings[0]
	month := subStrings[1]

	monthName := ""
	switch month {
	case "01":
		monthName = "Januari"
	case "02":
		monthName = "Februari"
	case "03":
		monthName = "Maret"
	case "04":
		monthName = "April"
	case "05":
		monthName = "Mei"
	case "06":
		monthName = "Juni"
	case "07":
		monthName = "Juli"
	case "08":
		monthName = "Agustus"
	case "09":
		monthName = "September"
	case "10":
		monthName = "Oktober"
	case "11":
		monthName = "November"
	case "12":
		monthName = "Desember"
	}

	formattedDate := fmt.Sprintf("%s %s", monthName, year)
	return formattedDate
}
