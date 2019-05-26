package jecnaapi

import (
	"testing"
)

func TestStringRankDo(t *testing.T) {
	rank := stringRank{
		strings: []string{"Lorem ipsum dolor sir amet", "Suplování 2018-19", "Suplování 2017-18"},
		rules: []stringRankRule{
			{stringHas: "Suplování", scoreIfTrue: 5, scoreIfFalse: -5},
			{stringHas: "2018", scoreIfTrue: 20},
			{stringHas: "18", scoreIfTrue: 100},
			{stringHas: "17", scoreIfTrue: -100},
			{stringHas: "2017", scoreIfTrue: -100},
			{stringHas: "2019", scoreIfTrue: 20},
		},
		continuous: true,
	}

	result := rank.do()

	for _, v := range result {
		var expectedScore int64
		switch v.string {
		case "Lorem ipsum dolor sir amet":
			expectedScore = -5
		case "Suplování 2018-19":
			expectedScore = 25
		case "Suplování 2017-18":
			expectedScore = 5
		}
		if v.score != expectedScore {
			t.Errorf("Score u stringu '%v' se neshoduje: ocekavano %+v, ma %+v", v.string, expectedScore, v.score)
		}
	}

	rank.continuous = false
	result = rank.do()

	for _, v := range result {
		var expectedScore int64
		switch v.string {
		case "Lorem ipsum dolor sir amet":
			expectedScore = -5
		case "Suplování 2018-19":
			expectedScore = 125
		case "Suplování 2017-18":
			expectedScore = -95
		}
		if v.score != expectedScore {
			t.Errorf("Score u stringu '%v' se neshoduje: ocekavano %+v, ma %+v", v.string, expectedScore, v.score)
		}
	}
}
