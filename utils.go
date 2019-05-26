package jecnaapi

import (
	"strings"
	"time"
)

var timeNow = time.Now

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func mesicCesky(mesic int) string {
	switch mesic {
	case 1:
		return "leden"
	case 2:
		return "únor"
	case 3:
		return "březen"
	case 4:
		return "duben"
	case 5:
		return "květen"
	case 6:
		return "červen"
	case 7:
		return "červenec"
	case 8:
		return "srpen"
	case 9:
		return "září"
	case 10:
		return "říjen"
	case 11:
		return "listopad"
	case 12:
		return "prosinec"
	default:
		return "chybec"
	}
}

func zacatekTohotoSkolnihoRoku(loc *time.Location) int {
	if int(time.Now().In(loc).Month()) >= 9 {
		return timeNow().In(loc).Year()
	} else {
		return timeNow().In(loc).Year() - 1
	}
}

func MakeSuplovaniTextReadable(s string) string {
	return strings.Replace(strings.Replace(strings.TrimSpace(strings.Replace(s, "\n", " ", -1)), "   ", " - ", -1), "  ", " - ", -1)
}
