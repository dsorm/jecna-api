// Podpurna miniknihovnicka pro vnitrni potreby Jecna API
// ktera nam usnadnuje hledani spravnych slozek pro ziskani
// momentalniho suplarchu, protoze zjevne mit sjednoceny
// format nazvu je prilis tezke pro Jecnou

package jecnaapi

import (
	"strings"
)

type stringRank struct {
	strings    []string
	rules      []stringRankRule
	continuous bool // pokud je true, postupne orezava string o dane stringy, kdezto pokud je false, tak pokazde orezava o puvodni string
}

type stringRankResult []stringRankResultRow

type stringRankRule struct {
	stringHas    string // kdyz string obsahuje tento string
	scoreIfTrue  int64  // tak za to dostane tento pocet bodu (lze i zaporny ci nula)
	scoreIfFalse int64  // a pokud nebsahuje, tak dostane tento pocet bodu
}

type stringRankResultRow struct {
	string string
	score  int64
}

func (sr *stringRank) do() stringRankResult {
	results := make(stringRankResult, len(sr.strings))
	var substr string
	var score int64

	if sr.continuous {
		var lastsubstr string
		for k, v := range sr.strings {
			lastsubstr = v
			substr = v
			score = 0
			for _, v2 := range sr.rules {
				substr = strings.Replace(lastsubstr, v2.stringHas, "", -1)
				if len(substr) < len(lastsubstr) {
					score += v2.scoreIfTrue
				} else {
					score += v2.scoreIfFalse
				}
				lastsubstr = substr
			}
			results[k] = stringRankResultRow{string: v, score: score}
		}
	} else {
		for k, v := range sr.strings {
			score = 0
			for _, v2 := range sr.rules {
				substr = strings.Replace(v, v2.stringHas, "", -1)
				if len(substr) < len(v) {
					score += v2.scoreIfTrue
				} else {
					score += v2.scoreIfFalse
				}
			}
			results[k] = stringRankResultRow{string: v, score: score}
		}
	}
	return results
}

/*
func (sr stringRankResult) HighestScore() stringRankResultRow {
	var highestscore int64 = -1
	highestscoreresult := -1
	for k, v := range sr {
		if v.score > highestscore {
			highestscoreresult = k
		}
	}
	return sr[highestscoreresult]
}

func (sr stringRankResult) LowestScore() stringRankResultRow {
	var lowestscore int64 = math.MaxInt64
	highestscoreresult := -1
	for k, v := range sr {
		if v.score > lowestscore {
			highestscoreresult = k
		}
	}
	return sr[highestscoreresult]
}
*/
