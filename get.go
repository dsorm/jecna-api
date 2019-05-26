package jecnaapi

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/extrame/xls"
	"github.com/imroc/req"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func (jecnak *Jecnak) GetSuplarchContent(suplarch Suplarch) SuplarchObsah {
	/*
		client := &http.Client{}
		request, err := http.NewRequest("GET", suplarch.url, nil)
		handleErr(err)
		request.Header.Add("Cookie", jecnak.LoginSession)
		resp, err := client.Do(request)
		handleErr(err)
	*/

	resp := jecnak.getHttpResponse(suplarch.url)
	// ioutil.ReadAll() prevede io.ReadSeeker na []byte, a bytes.NewReader() z toho udela io.Reader
	suplarchByte, err := ioutil.ReadAll(resp.Body)
	handleErr(err)
	suplarchXls, err := xls.OpenReader(bytes.NewReader(suplarchByte), "utf-8")
	handleErr(err)
	handleErr(resp.Body.Close())
	// vybrat prvni arch z xls souboru, stejne jich tam moc neni
	sheet := suplarchXls.GetSheet(0)

	tridy := make([]string, 0, sheet.MaxRow)
	suplRowStart := 0
	//suplRowEnd := 0
	var tridyColumn int
	if sheet.Row(8).Col(0) == "" {
		tridyColumn = 1
	} else {
		tridyColumn = 0
	}
	suplColumnStart := 1 + tridyColumn
	suplColumnEnd := 10
	for i := 0; i < int(sheet.MaxRow); i++ {
		if sheet.Row(i).Col(tridyColumn+1) == "1" && (len(tridy) > 0) {
			suplColumnContent := "necozatim"
			i2 := 3
			for !(suplColumnContent == "") {
				suplColumnContent = sheet.Row(i).Col(i2)
				i2++
			}
			suplColumnEnd = i2 - 2
			//suplRowEnd = i -1
			break
		}
		if !strings.Contains(sheet.Row(i).Col(tridyColumn), "/") {
			suplRowStart = i + 1
			continue
		}

		tridy = append(tridy, strings.TrimSpace(strings.Split(sheet.Row(i).Col(tridyColumn), "/")[0]))
	}

	suplarchObsah := make(SuplarchObsah, 0, len(tridy))

	for k, v := range tridy {
		suplyString := make([]string, 0, suplColumnEnd-suplColumnStart)

		for i := suplColumnStart; i <= suplColumnEnd; i++ {
			suplyString = append(suplyString, sheet.Row(k+suplRowStart).Col(i))
		}

		suplarchObsah = append(suplarchObsah, SuplarchTrida{
			trida: v,
			supl:  suplyString,
		})
	}
	return suplarchObsah
}

func (jecnak *Jecnak) GetSuplarchSeznam() SuplarchSeznam {
	loc, _ := time.LoadLocation("CET") // CET - Central European Time
	var suplarchy SuplarchSeznam
	currentMonth := mesicCesky(int(time.Now().In(loc).Month()))
	monthWeekAgo := mesicCesky(int(time.Now().In(loc).AddDate(0, 0, -7).Month()))
	lastSchoolYear := zacatekTohotoSkolnihoRoku(loc)
	currentSchoolYear := lastSchoolYear + 1
	currentYear := time.Now().In(loc).Year()
	yearWeekAgo := time.Now().In(loc).AddDate(0, 0, -7).Year()
	/*
		fmt.Println("Current month:", currentMonth)
		fmt.Println("Month week ago:", monthWeekAgo)
		fmt.Println("Current year:", strconv.Itoa(currentYear))
		fmt.Println("Year week ago", strconv.Itoa(yearWeekAgo))
		fmt.Println("Current school year:", strconv.Itoa(lastSchoolYear)+"/"+strconv.Itoa(currentSchoolYear))
	*/
	response, err := req.Get(siteURL+"dokumenty/skola/pro-studenty/", req.Header{
		"User-Agent": "Mozilla/5.0",
		"Cookie":     jecnak.LoginSession,
	})
	handleErr(err)
	document, err := goquery.NewDocumentFromReader(response.Response().Body)
	handleErr(err)

	// hledame slozku se skolnim rokem na webu
	highestScore := 0
	highestScoreURL := ""
	document.Find("ul.documents a.dir").Each(func(i int, selection *goquery.Selection) {
		puvodniText := strings.Replace(strings.ToLower(selection.Find("span.label").Text()), "&nbsp;", "", -1)
		sr := stringRank{
			strings: []string{puvodniText},
			rules: []stringRankRule{
				{stringHas: strconv.Itoa(currentSchoolYear), scoreIfTrue: 15},           // 2019
				{stringHas: strconv.Itoa(lastSchoolYear), scoreIfTrue: 15},              // 2018
				{stringHas: strconv.Itoa(lastSchoolYear - 1), scoreIfTrue: -10},         // 2017
				{stringHas: strconv.Itoa(currentSchoolYear % 100), scoreIfTrue: 15},     // 19
				{stringHas: strconv.Itoa(lastSchoolYear % 100), scoreIfTrue: 15},        // 18
				{stringHas: strconv.Itoa((lastSchoolYear - 1) % 100), scoreIfTrue: -10}, // 17
				{stringHas: "suplování a akce", scoreIfTrue: 10},
				{stringHas: "suplování", scoreIfTrue: 7},
			},
			continuous: true,
		}
		result := sr.do()

		if int(result[0].score) > highestScore {
			url, _ := selection.Attr("href")
			highestScore = int(result[0].score)
			highestScoreURL = strings.TrimPrefix(url, "/")
			// fmt.Printf("Nalezena nova nejpravdepodobnejsi slozka: '%v', s pravdepodobnosti %+v (%v)\n", highestScoreText, highestScore, highestScoreURL)
		}
	})
	currentYearURL := siteURL + highestScoreURL

	// hledame slozku se spravnym mesicem ted
	response, err = req.Get(currentYearURL, req.Header{
		"User-Agent": "Mozilla/5.0",
		"Cookie":     jecnak.LoginSession,
	})
	handleErr(err)
	document, err = goquery.NewDocumentFromReader(response.Response().Body)
	handleErr(err)
	// fmt.Println(response.Response().Status, response.Request().URL.String())
	highestScore = 0
	highestScore = 0
	var currentMonthURL string
	document.Find("ul.documents a.dir").Each(func(i int, selection *goquery.Selection) {
		puvodniText := strings.Replace(strings.ToLower(selection.Find("span.label").Text()), "&nbsp;", "", -1)

		sr := stringRank{
			strings: []string{puvodniText},
			rules: []stringRankRule{
				{stringHas: currentMonth, scoreIfTrue: 5},                    // zkontrolovat, zda-li obsahuje tento mesic ("unor")
				{stringHas: strconv.Itoa(currentYear), scoreIfTrue: 5},       // zkontrolovat, zda-li obsahuje tento rok ("2019")
				{stringHas: strconv.Itoa(currentYear % 100), scoreIfTrue: 5}, // zkontrolovat, zda-li obsahuje tento rok ("19")
			},
			continuous: true,
		}

		results := sr.do()

		if int(results[0].score) > highestScore {
			url, _ := selection.Attr("href")
			highestScore = int(results[0].score)
			currentMonthURL = strings.TrimPrefix(url, "/")
			// fmt.Printf("Nalezena nova nejpravdepodobnejsi slozka s mesicem: '%v', s pravdepodobnosti %+v (%v)\n", highestScoreText, highestScore, currentMonthURL)
		}
	})

	// hledame jednotlive suplarchy
	if highestScore == 0 {
		currentMonthURL = highestScoreURL
	}
	response, err = req.Get(siteURL+currentMonthURL, req.Header{
		"User-Agent": "Mozilla/5.0",
		"Cookie":     jecnak.LoginSession,
	})
	handleErr(err)
	document, err = goquery.NewDocumentFromReader(response.Response().Body)
	handleErr(err)

	document.Find("ul.documents a.file").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if len(suplarchy) == 5 {
			return false
		}
		text := selection.Text()
		cutText := strings.TrimSpace(strings.TrimPrefix(text, "Suplarch"))
		url, _ := selection.Attr("href")
		url = siteURL + strings.TrimPrefix(url, "/")
		if len(cutText) < len(text) {
			// fmt.Printf("Nalezen nový suplarch: %v (%v)\n", cutText, url)
			suplarchy = append(suplarchy, Suplarch{name: cutText, url: url})
		}
		return true
	})

	if len(suplarchy) < 5 {
		// hledame slozku se spravnym mesicem pred tydnem
		response, err = req.Get(currentYearURL, req.Header{
			"User-Agent": "Mozilla/5.0",
			"Cookie":     jecnak.LoginSession,
		})
		handleErr(err)
		document, err = goquery.NewDocumentFromReader(response.Response().Body)
		handleErr(err)
		// fmt.Println(response.Response().Status, response.Request().URL.String())
		highestScore = 0
		highestScore = 0
		var lastMonthURL string
		document.Find("ul.documents a.dir").Each(func(i int, selection *goquery.Selection) {
			puvodniText := strings.Replace(strings.ToLower(selection.Find("span.label").Text()), "&nbsp;", "", -1)

			sr := stringRank{
				strings: []string{puvodniText},
				rules: []stringRankRule{
					{stringHas: monthWeekAgo, scoreIfTrue: 5},                    // zkontrolovat, zda-li obsahuje tento mesic ("leden")
					{stringHas: strconv.Itoa(yearWeekAgo), scoreIfTrue: 5},       // zkontrolovat, zda-li obsahuje tento rok ("2019")
					{stringHas: strconv.Itoa(yearWeekAgo % 100), scoreIfTrue: 5}, // zkontrolovat, zda-li obsahuje tento rok ("19")
				},
				continuous: true,
			}

			results := sr.do()

			if int(results[0].score) > highestScore {
				url, _ := selection.Attr("href")
				highestScore = int(results[0].score)
				lastMonthURL = strings.TrimPrefix(url, "/")
				//	fmt.Printf("Nalezena nova nejpravdepodobnejsi slozka s mesicem: '%v', s pravdepodobnosti %+v (%v)\n", highestScoreText, highestScore, currentMonthURL)
			}

		})

		// hledame suplarchy z mesice o tyden zpet
		response, err = req.Get(siteURL+lastMonthURL, req.Header{
			"User-Agent": "Mozilla/5.0",
			"Cookie":     jecnak.LoginSession,
		})
		handleErr(err)
		document, err = goquery.NewDocumentFromReader(response.Response().Body)
		handleErr(err)

		document.Find("ul.documents a.file").EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if len(suplarchy) == 5 {
				return false
			}
			text := selection.Text()
			cutText := strings.TrimSpace(strings.TrimPrefix(text, "Suplarch"))
			url, _ := selection.Attr("href")
			url = siteURL + strings.TrimPrefix(url, "/")
			if len(cutText) < len(text) {
				// fmt.Printf("Nalezen nový suplarch: %v (%v)\n", cutText, url)
				suplarchy = append(suplarchy, Suplarch{name: cutText, url: url})
			}
			return true
		})
	}

	return suplarchy
}

func (jecnak *Jecnak) GetPrichody() Prichody {
	document := jecnak.getGoqueryDocument(siteURL + "absence/passing-student")

	var po Prichody
	document.Find("tbody tr").Each(func(iteration int, subselection *goquery.Selection) {
		datum := strings.Split(subselection.Find("td.date").Text(), " (")[0] // odstrani napr " (Patek)" na konci, protoze je nam to k nicemu
		text := subselection.Find("td").Eq(1).Text()

		po = append(po, prichodyRow{
			Datum:         datum,
			PrichodOdchod: text,
		})
	})

	return po
}

func (jecnak *Jecnak) GetOmluvnyList() OmluvnyList {
	document := jecnak.getGoqueryDocument(siteURL + "absence/student")

	var hodiny OmluvnyList
	document.Find("tbody tr").Each(func(iteration int, subselection *goquery.Selection) {
		datum := strings.Split(subselection.Find("td.date").Text(), " (")[0] // odstrani napr " (Patek)" na konci, protoze je nam to k nicemu
		text := subselection.Find("td.count").Text()

		hodiny = append(hodiny, omluvnyListRow{
			Datum: strings.Trim(datum, "\n            "), // opravuje divny bug, ktery zpusobuje ze na konci a na zacatku je tenhle string
			Text:  text,
		})
	})

	return hodiny

}

func (jecnak *Jecnak) GetJecnakInfo() {
	document := jecnak.getGoqueryDocument(siteURL + "student/" + jecnak.LoginName)

	document.Find("table.userprofile tr").Each(func(iteration int, subselection *goquery.Selection) {
		subselectionText := subselection.Find("span.value").Text()
		if len(subselectionText) < 1 {
			// nektere pole jsou rozklikavaci a nejsou ve span.value
			subselectionText = subselection.Find("a.link span.label").Text()
		}
		switch subselection.Find("span.label").First().Text() {
		case "Celé jméno":
			jecnak.CeleJmeno = subselectionText
		case "Uživatelské jméno":
			jecnak.LoginName = subselectionText
		case "Věk":
			vek, _ := strconv.Atoi(strings.Replace(subselectionText, " let", "", 1))
			jecnak.Vek = uint8(vek)
		case "Narození":
			split := strings.Split(subselectionText, ", ")
			jecnak.MistoNarozeni = split[1]
			jecnak.Narozeniny = split[0]
		case "Telefon":
			telefon, _ := strconv.Atoi(strings.Replace(subselectionText, " ", "", -1))
			jecnak.Telefon = uint32(telefon)
		case "Trvlá adresa", "Trvalá adresa": // ne, toto není překlep, aspoň z mé strany
			split := strings.Split(subselectionText, ", ")
			jecnak.TrvalaAdresaUlice = split[0]
			jecnak.TrvalaAdresaMesto = split[1]
			jecnak.TrvalaAdresaPSC = split[2]
		case "Třída, skupiny":
			jecnak.SkupinyTrida = subselectionText
		case "Číslo v tříd. výkazu", "Číslo v třídním výkazu":
			cislo, _ := strconv.Atoi(subselectionText)
			jecnak.CisloVTridnimVykazu = uint8(cislo)
		case "Soukromý e-mail":
			jecnak.SoukromyEmail = subselectionText
		case "Školní e-mail":
			jecnak.SkolniEmail = strings.Replace(subselectionText, " "+subselection.Find("span.value em").Text(), "", 1) // odebrat " (přeposílán na soukromý)"
		}
	})
}

func (jecnak *Jecnak) GetRozvrh() Rozvrh {
	document := jecnak.getGoqueryDocument(siteURL + "timetable/class")

	var timetable Rozvrh
	document.Find("table.timetable tr").Each(func(iteration int, subselection *goquery.Selection) {
		if iteration == 0 {
			subselection.Find("th.period").Each(func(subiteration int, subsubselection *goquery.Selection) {
				// "subsubselection.Clone().Children().Remove().End().Text()" vrati pouze Text daneho prvku bez jeho deti
				timetable.Casy = append(timetable.Casy, RozvrhCas{CisloHodiny: subsubselection.Clone().Children().Remove().End().Text(), Doba: subsubselection.Find("span.time").Text()})
			})
		} else {
			var rozvrhden []RozvrhHodina
			subselection.Find("td").Each(func(subiteration int, subsubselection *goquery.Selection) {
				rozvrhden = append(rozvrhden, RozvrhHodina{ZkratkaHodiny: subsubselection.Find("span.subject").Text()})
			})
			switch iteration {
			case 1:
				timetable.Pondeli = rozvrhden
			case 2:
				timetable.Utery = rozvrhden
			case 3:
				timetable.Streda = rozvrhden
			case 4:
				timetable.Ctvrtek = rozvrhden
			case 5:
				timetable.Patek = rozvrhden
			}
		}
	})
	return timetable
}

func (jecnak *Jecnak) GetZnamky() Znamky {
	document := jecnak.getGoqueryDocument(siteURL + "score/student")

	scores := Znamky{}
	document.Find("table.score tbody tr").Each(func(iteration int, subselection *goquery.Selection) {
		text := subselection.Find("th").Text()

		var split []string
		var zkratka, nazevPredmetu string
		var subjectScores []znamka
		var final uint8 = 0

		if text != "Chování" { // chování nemá žádnou zkratku - zpusobilo by panic
			split = strings.Split(text, " (")
			zkratka = strings.Replace(split[1], ")", "", 1)
			nazevPredmetu = split[0]
		} else {
			nazevPredmetu = text
			zkratka = ""
		}

		subselection.Find("a.score").Each(func(iteration int, subsubselection *goquery.Selection) {
			if subsubselection.HasClass("scoreFinal") { // ta znamka uplne napravo
				preconvert, _ := strconv.Atoi(subsubselection.Text())
				final = uint8(preconvert)
			} else {
				score := znamka{}
				preconvert, _ := strconv.Atoi(subsubselection.Find("span.value").Text())

				presplit, _ := subsubselection.Attr("title")
				titleSplit := strings.Split(presplit, "(")

				score.Hodnota = uint8(preconvert)
				score.Vyucujici = subsubselection.Find("span.employee").Text()
				score.Datum = strings.Replace(strings.Split(titleSplit[1], ", ")[0], "(", "", -1)
				score.Mala = subsubselection.HasClass("scoreSmall")
				score.Popis = strings.TrimSpace(titleSplit[0])
				subjectScores = append(subjectScores, score)
				final = 100
			}
		})

		scores = append(scores, znamkyPredmet{
			Nazev:   nazevPredmetu,
			Zkratka: zkratka,
			Znamky:  subjectScores,
			Final:   final,
		})

	})

	return scores
}

func (jecnak *Jecnak) GetJecnitel(jecnitelovaZkratka string) Jecnitel {
	document := jecnak.getGoqueryDocument(siteURL + "ucitel/" + jecnitelovaZkratka)

	ucitel := Jecnitel{}
	document.Find("table.userprofile tr").Each(func(iteration int, subselection *goquery.Selection) {
		subselectionText := subselection.Find("span.value").Text()
		if len(subselectionText) < 1 {
			// nektere pole jsou rozklikavaci a nejsou ve span.value
			subselectionText = subselection.Find("span.label").Eq(1).Text() // .eq(1) vybere druhy span.label element (array-like cislovani)
		}
		switch subselection.Find("span.label").First().Text() { // first protoze nektere udaje jsou a.link, ktery je obaleny v span.label, tudiz by sebralo i hodnotu pole jakozto Nazev (napr. "E-mailnovak@spsejecna.cz")
		case "Jméno":
			ucitel.Jmeno = subselectionText
		case "Zkratka":
			ucitel.Zkratka = subselectionText
		case "Uživatelské jméno":
			ucitel.UzivatelskeJmeno = subselectionText
		case "E-mail":
			ucitel.Email = subselectionText
		case "Telefon":
			ucitel.Telefon = strings.Replace(strings.TrimSpace(subselectionText), "  ", " ", -1) // orizneme mezeru na zacatku, a vymenime vsechny dvojmezery za mezery
		case "Kabinet":
			ucitel.Kabinet = subselectionText
		case "Konzultační hodiny":
			ucitel.KonzultacniHodiny = subselectionText
		case "Soukromý e-mail":
			ucitel.SoukromyEmail = subselectionText
		case "Třídní učitel":
			ucitel.TridniUcitel = subselectionText
		}
	})
	if jecnak.LoginValid && ucitel.Jmeno != "" {
		// fotky ucitelu vidi jen prihlaseni
		ucitel.ImgUrl, _ = document.Find("div.image img").Attr("src")
	}
	return ucitel
}

func (jecnak *Jecnak) GetSeznamJecnitelu() SeznamJecnitelu {
	document := jecnak.getGoqueryDocument(siteURL + "ucitel")

	var jecnitele []SeznamJecniteluRow
	document.Find("div.contentLeftColumn ul li").Each(func(iteration int, subselection *goquery.Selection) {
		attribut, _ := subselection.Find("a").Attr("href")
		attribut = strings.Title(strings.ToLower(strings.Replace(attribut, "/ucitel/", "", -1))) // odebere /ucitel/, udela maly pismena, a potom na zacatku da velke pismeno
		jecnitele = append(jecnitele, SeznamJecniteluRow{subselection.Text(), attribut})
	})
	document.Find("div.contentRightColumn ul li").Each(func(iteration int, subselection *goquery.Selection) {
		attribut, _ := subselection.Find("a").Attr("href")
		attribut = strings.Title(strings.ToLower(strings.Replace(attribut, "/ucitel/", "", -1))) // odebere /ucitel/, udela maly pismena, a potom na zacatku da velke pismeno
		jecnitele = append(jecnitele, SeznamJecniteluRow{subselection.Text(), attribut})
	})
	return jecnitele
}
