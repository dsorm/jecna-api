package jecnaapi

import (
	"testing"
	"time"
)

type chyba struct{}

func (chyba) Error() string {
	return "chyba"
}

func TestHandleErr(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Nebyl zadny panic")
			}
		}()
		handleErr(chyba{})
	}()
}

func TestMesicCesky(t *testing.T) {
	for i := 1; i < 14; i++ {
		output := mesicCesky(i)

		var expected string
		switch i {
		case 1:
			expected = "leden"
		case 2:
			expected = "únor"
		case 3:
			expected = "březen"
		case 4:
			expected = "duben"
		case 5:
			expected = "květen"
		case 6:
			expected = "červen"
		case 7:
			expected = "červenec"
		case 8:
			expected = "srpen"
		case 9:
			expected = "září"
		case 10:
			expected = "říjen"
		case 11:
			expected = "listopad"
		case 12:
			expected = "prosinec"
		case 13:
			expected = "chybec"
		}

		if !(output == expected) {
			t.Errorf("Ocekavano '%v', ale vysledek je '%v' (vstup: %+v)", expected, output, i)
		}
	}
}

func TestZacatekTohotoSkolnihoRoku(t *testing.T) {
	loc, _ := time.LoadLocation("Europe/Prague")
	output := zacatekTohotoSkolnihoRoku(loc)
	/*
		timeNow, _ = time.Parse("2.1.2006", "10.9.2018")
		if got := zacatekTohotoSkolnihoRoku(loc); got != 2018 {
			t.Errorf("Ocekavano 2018, bylo vraceno %v", got)
		}

		timeNow, _ = time.
		if got := zacatekTohotoSkolnihoRoku(loc); got != 2016 {
			t.Errorf("Ocekavano 2016, bylo vraceno %v", got)
		}
	*/
	if int(time.Now().In(loc).Month()) >= 9 {
		if output != time.Now().In(loc).Year() {
			t.Fail()
		}
	} else {
		if output != time.Now().In(loc).Year()-1 {
			t.Fail()
		}
	}
}

func TestMakeSuplovaniTextReadable(t *testing.T) {
	const str string = "\n  "
	if len(MakeSuplovaniTextReadable(str)) != 0 {
		t.Error()
	}
}
