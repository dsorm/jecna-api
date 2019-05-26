package jecnaapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var mockJecnak = Jecnak{LoginName: "novak", LoginPass: "123", LoginSession: "JSESSIONID=nejakejsession", LoginValid: true}

type httpGetRoute struct {
	url  string
	file string
}

type httpGetRoutes []httpGetRoute

func testGetRequestMultiRoute(t *testing.T, jecnak Jecnak, routes httpGetRoutes, getStructFunc func() interface{}, expectedStruct interface{}) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("cookie") != jecnak.LoginSession+"; role=student;" {
			t.Errorf("Cookie se neshoduje, ocekavano '%v', je '%v'", jecnak.LoginSession+"; role=student;", req.Header.Get("cookie"))
		}

		validRoute := false
		for _, v := range routes {
			if req.URL.Path == "/"+v.url {
				validRoute = true
				if body, err := ioutil.ReadFile(v.file); err != nil {
					t.Errorf("Chyba testu: nelze otevrit soubor '%v'\n", v.file)
				} else {
					if !strings.HasSuffix(strings.ToLower(v.file), ".html-test") {
						Openfile, err := os.Open(v.file)
						if err != nil {
							t.Errorf("Chyba testu: soubor nebyl nalezen: '%v'", v.file)
							return
						}
						FileHeader := make([]byte, 512)
						_, err = Openfile.Read(FileHeader)
						handleErr(err)
						FileContentType := http.DetectContentType(FileHeader)
						FileStat, err := Openfile.Stat()
						handleErr(err)
						FileSize := strconv.FormatInt(FileStat.Size(), 10)
						rw.Header().Set("Content-Disposition", "attachment; filename="+v.file)
						rw.Header().Set("Content-Type", FileContentType)
						rw.Header().Set("Content-Length", FileSize)
						handleErr(Openfile.Close()) // Close after function return
					}
					if _, err := rw.Write(body); err != nil {
						t.Errorf("Chyba testu: chyba pri psani odpovedi HTTP\n")
					}
				}

			}
		}

		if !validRoute {
			t.Errorf("Nebyl nalezen route s touto url: %v\n", req.URL.Path)
		}
	}))
	siteURL = server.URL + "/"
	defer server.Close()
	gotStruct := getStructFunc()

	if !reflect.DeepEqual(gotStruct, expectedStruct) {
		typ := reflect.TypeOf(gotStruct)
		expectedStructProcessed := strings.Replace(strings.Trim(strings.Replace(fmt.Sprintf("%#v", expectedStruct), "jecnaapi.", "", -1), `"`), `\\"`, `"`, -1)
		gotStructProcessed := strings.Replace(strings.Trim(strings.Replace(fmt.Sprintf("%#v", gotStruct), "jecnaapi.", "", -1), `"`), `\`+`"`, `"`, -1)
		t.Errorf("\nNebyly vraceny ocekavane udaje u typu %v{}\n Ocekavano: %v\n Bylo vraceno: %v", typ.Name(), expectedStructProcessed, gotStructProcessed)
	}
}

func testGetRequest(t *testing.T, jecnak Jecnak, expectedURL string, htmlFile string, getStructFunc func() interface{}, expectedStruct interface{}) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("cookie") != jecnak.LoginSession+"; role=student;" {
			t.Errorf("Cookie se neshoduje, ocekavano '%v', je '%v'", jecnak.LoginSession+"; role=student;", req.Header.Get("cookie"))
		}
		if req.URL.Path != "/"+expectedURL {
			t.Errorf("URL se neshoduje, ocekavano '%v', je '%v'", "/"+expectedURL, req.URL.Path)
		}
		body, _ := ioutil.ReadFile(htmlFile)
		_, _ = rw.Write(body)
	}))
	siteURL = server.URL + "/"
	defer server.Close()

	gotStruct := getStructFunc()
	equal := reflect.DeepEqual(gotStruct, expectedStruct)

	if !equal {
		typ := reflect.TypeOf(gotStruct)
		expectedStructProcessed := strings.Replace(strings.Trim(strings.Replace(fmt.Sprintf("%#v", expectedStruct), "jecnaapi.", "", -1), `"`), `\\"`, `"`, -1)
		gotStructProcessed := strings.Replace(strings.Trim(strings.Replace(fmt.Sprintf("%#v", gotStruct), "jecnaapi.", "", -1), `"`), `\`+`"`, `"`, -1)
		t.Errorf("\nNebyly vraceny ocekavane udaje u typu %v{}\n Ocekavano: %v\n Bylo vraceno: %v", typ.Name(), expectedStructProcessed, gotStructProcessed)
		// t.Errorf("\n Nezpracovany vystup: %#v", gotStruct)
	}

}

func TestJecnak_GetSuplarchContent(t *testing.T) {
	expected := SuplarchObsah{
		SuplarchTrida{trida: "A1", supl: []string{"", "", "", "", "\n", "82 Braunova\n úklid uč", "", "", ""}},
		SuplarchTrida{trida: "C1a", supl: []string{"", "", "", "", "", "1/2  N 80\n Hi", "1/2  N 80\n Hi", "", ""}},
		SuplarchTrida{trida: "C1b", supl: []string{"IT 70\n St ", "IT 70\nSt (Le) + ", "IT 70\nSt (Le) + ", "TVy 90\nLe (Braunova) +", "1/2 A 75\nLe (St) +\n 1/2 A 75\nMo +", "", "", "", ""}},
		SuplarchTrida{trida: "E1a", supl: []string{"", "", "", "", "", "60 Ei\n úklid uč", "", "", ""}},
		SuplarchTrida{trida: "E1b", supl: []string{"", "", "", "M 50 +\nHi (Mo)", "F 50 +\nHi (Hm)", " 50 Mo", "50 St", "", ""}},
	}
	routes := httpGetRoutes{{
		url:  "SPSE-Jecna_Suplarch+so_21_09_02.XLS",
		file: "_test/SPSE-Jecna_Suplarch so_21_09_02.XLS",
	}}

	getFunc := func() interface{} {
		suplarch := Suplarch{
			name: "so_21_09_02",
			url:  siteURL + "SPSE-Jecna_Suplarch+so_21_09_02.XLS",
		}
		return mockJecnak.GetSuplarchContent(suplarch)
	}

	testGetRequestMultiRoute(t, mockJecnak, routes, getFunc, expected)
}

func TestJecnak_GetPrichody(t *testing.T) {
	expected := Prichody{
		prichodyRow{Datum: "1.2.\u00a0(pátek)", PrichodOdchod: ""},
		prichodyRow{Datum: "4.2.\u00a0(pondělí)", PrichodOdchod: ""},
		prichodyRow{Datum: "5.2.\u00a0(úterý)", PrichodOdchod: "Odchod 14:26"},
		prichodyRow{Datum: "6.2.\u00a0(středa)", PrichodOdchod: "Odchod 14:12"},
		prichodyRow{Datum: "7.2.\u00a0(čtvrtek)", PrichodOdchod: "Příchod 7:26, Odchod 14:12"},
		prichodyRow{Datum: "8.2.\u00a0(pátek)", PrichodOdchod: "Odchod 13:21"},
		prichodyRow{Datum: "11.2.\u00a0(pondělí)", PrichodOdchod: ""},
		prichodyRow{Datum: "12.2.\u00a0(úterý)", PrichodOdchod: ""},
		prichodyRow{Datum: "13.2.\u00a0(středa)", PrichodOdchod: ""},
		prichodyRow{Datum: "14.2.\u00a0(čtvrtek)", PrichodOdchod: ""},
		prichodyRow{Datum: "15.2.\u00a0(pátek)", PrichodOdchod: ""},
		prichodyRow{Datum: "18.2.\u00a0(pondělí)", PrichodOdchod: ""},
		prichodyRow{Datum: "19.2.\u00a0(úterý)", PrichodOdchod: ""},
		prichodyRow{Datum: "20.2.\u00a0(středa)", PrichodOdchod: ""},
		prichodyRow{Datum: "21.2.\u00a0(čtvrtek)", PrichodOdchod: ""},
		prichodyRow{Datum: "22.2.\u00a0(pátek)", PrichodOdchod: ""},
		prichodyRow{Datum: "25.2.\u00a0(pondělí)", PrichodOdchod: ""},
		prichodyRow{Datum: "26.2.\u00a0(úterý)", PrichodOdchod: ""},
		prichodyRow{Datum: "27.2.\u00a0(středa)", PrichodOdchod: ""},
		prichodyRow{Datum: "28.2.\u00a0(čtvrtek)", PrichodOdchod: ""},
	}
	testGetRequest(t, mockJecnak, "absence/passing-student",
		"_test/prichody.html-test", func() interface{} { return mockJecnak.GetPrichody() }, expected)

	routes := httpGetRoutes{{
		url:  "absence/passing-student",
		file: "_test/prichody.html-test",
	}}

	testGetRequestMultiRoute(t, mockJecnak, routes, func() interface{} { return mockJecnak.GetPrichody() }, expected)
}

func TestJecnak_GetOmluvnyList(t *testing.T) {
	expected := OmluvnyList{
		omluvnyListRow{Datum: "29.1.", Text: "7 hodin"},
		omluvnyListRow{Datum: "22.1.", Text: "1 hodina"},
		omluvnyListRow{Datum: "16.1.", Text: "3 hodiny"},
		omluvnyListRow{Datum: "11.1.", Text: "4 hodiny z toho 4 neomluvené"},
		omluvnyListRow{Datum: "10.1.", Text: "7 hodin z toho 7 neomluvených"},
		omluvnyListRow{Datum: "4.1.", Text: "2 hodiny z toho 2 neomluvené"},
		omluvnyListRow{Datum: "13.12.", Text: "1 pozdní příchod z toho 1 neomluvený"},
		omluvnyListRow{Datum: "4.12.", Text: "7 hodin"},
	}
	testGetRequest(t, mockJecnak, "absence/student",
		"_test/omluvnylist.html-test", func() interface{} { return mockJecnak.GetOmluvnyList() }, expected)
}

func TestJecnak_GetJecnakInfo(t *testing.T) {
	expected := Jecnak{
		LoginName:    "novak",
		LoginPass:    "123",
		LoginSession: "JSESSIONID=nejakejsession",
		LoginValid:   true, CeleJmeno: "Jan Novák",
		Vek:                 0x39,
		Narozeniny:          "1. 9. 1939",
		MistoNarozeni:       "Polsko/Německo",
		Telefon:             0x75bcd15,
		TrvalaAdresaUlice:   "Nekde 1234/5",
		TrvalaAdresaMesto:   "Praha 666",
		TrvalaAdresaPSC:     "12345",
		SkupinyTrida:        "E1a, skupiny: A2",
		CisloVTridnimVykazu: 0x9a,
		SoukromyEmail:       "novak@example.com",
		SkolniEmail:         "novak@spsejecna.example.com",
	}
	testGetRequest(t, mockJecnak, "student/novak",
		"_test/jecnak.html-test", func() interface{} { mockJecnak.GetJecnakInfo(); return mockJecnak }, expected)
}

func TestJecnak_GetRozvrh(t *testing.T) {
	expected := Rozvrh{
		Casy: []RozvrhCas{{CisloHodiny: "1", Doba: "7:30 - 8:15"},
			{CisloHodiny: "2", Doba: "8:25 - 9:10"},
			{CisloHodiny: "3", Doba: "9:20 - 10:05"},
			{CisloHodiny: "4", Doba: "10:20 - 11:05"},
			{CisloHodiny: "5", Doba: "11:15 - 12:00"},
			{CisloHodiny: "6", Doba: "12:10 - 12:55"},
			{CisloHodiny: "7", Doba: "13:05 - 13:50"},
			{CisloHodiny: "8", Doba: "14:00 - 14:45"},
			{CisloHodiny: "9", Doba: "14:55 - 15:40"},
			{CisloHodiny: "10", Doba: "15:50 - 16:35"}},

		Pondeli: []RozvrhHodina{{UcitelZkratka: "", ZkratkaHodiny: "AA", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "C", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "TV", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "ON", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "PAA", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "M", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "ZE", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""}},

		Utery: []RozvrhHodina{{UcitelZkratka: "", ZkratkaHodiny: "PA", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "D", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "M", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "ITIT", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "ITIT", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "AA", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "ZE", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""}},

		Streda: []RozvrhHodina{{UcitelZkratka: "", ZkratkaHodiny: "Ch", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "M", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "C", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "IT", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "PAA", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "ON", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "PA", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""}},

		Ctvrtek: []RozvrhHodina{{UcitelZkratka: "", ZkratkaHodiny: "Ch", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "TV", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "F", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "M", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "DCDCDC", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "DCDCDC", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "DCDCDC", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""}},

		Patek: []RozvrhHodina{{UcitelZkratka: "", ZkratkaHodiny: "ON", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "AA", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "F", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "C", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "D", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "M", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""},
			{UcitelZkratka: "", ZkratkaHodiny: "", Ucebna: "", Trida: "", Skupina: ""}},
	}
	testGetRequest(t, mockJecnak, "timetable/class",
		"_test/rozvrh.html-test", func() interface{} { return mockJecnak.GetRozvrh() }, expected)
}

func TestJecnak_GetZnamky(t *testing.T) {
	expected := Znamky{
		znamkyPredmet{Nazev: "Český jazyk a literatura", Zkratka: "C", Znamky: []znamka{
			{Popis: "", Vyucujici: "SC", Hodnota: 0x3, Datum: "13.09.2018", Mala: false},
			{Popis: "", Vyucujici: "SC", Hodnota: 0x3, Datum: "04.10.2018", Mala: false},
			{Popis: "diktát", Vyucujici: "SU", Hodnota: 0x3, Datum: "07.11.2018", Mala: false},
			{Popis: "diktát", Vyucujici: "SU", Hodnota: 0x2, Datum: "doplňkový) ", Mala: true},
			{Popis: "", Vyucujici: "SU", Hodnota: 0x0, Datum: "03.12.2018", Mala: false},
			{Popis: "zkoušení", Vyucujici: "SU", Hodnota: 0x2, Datum: "14.12.2018", Mala: false},
			{Popis: "diktát", Vyucujici: "SU", Hodnota: 0x2, Datum: "20.01.2019", Mala: false},
			{Popis: "test - literatura", Vyucujici: "SU", Hodnota: 0x2, Datum: "20.01.2019", Mala: false}},
			Final: 0x4},

		znamkyPredmet{Nazev: "Anglický jazyk", Zkratka: "A", Znamky: []znamka{
			{Popis: "Test 1: hodiny, předložky časové, přítomné časy", Vyucujici: "KO", Hodnota: 0x2, Datum: "01.10.2018", Mala: false},
			{Popis: "Minulý čas prostý", Vyucujici: "KO", Hodnota: 0x2, Datum: "16.10.2018", Mala: false},
			{Popis: "Slovicka Šikana a lekce 1", Vyucujici: "KO", Hodnota: 0x1, Datum: "05.11.2018", Mala: true},
			{Popis: "test slovesné vazby", Vyucujici: "KO", Hodnota: 0x0, Datum: "19.11.2018", Mala: true},
			{Popis: "test Unit 2- první část", Vyucujici: "KO", Hodnota: 0x1, Datum: "04.12.2018", Mala: true},
			{Popis: "Test Unit 2", Vyucujici: "KO", Hodnota: 0x1, Datum: "18.12.2018", Mala: false},
			{Popis: "domácí úkoly-1.pololetí", Vyucujici: "KO", Hodnota: 0x0, Datum: "14.01.2019", Mala: false}},
			Final: 0x2},

		znamkyPredmet{Nazev: "Dějepis", Zkratka: "D", Znamky: []znamka{
			{Popis: "", Vyucujici: "SA", Hodnota: 0x2, Datum: "20.10.2018", Mala: false},
			{Popis: "", Vyucujici: "SA", Hodnota: 0x2, Datum: "20.10.2018", Mala: false},
			{Popis: "", Vyucujici: "SA", Hodnota: 0x2, Datum: "06.11.2018", Mala: false},
			{Popis: "", Vyucujici: "SA", Hodnota: 0x4, Datum: "13.11.2018", Mala: false},
			{Popis: "", Vyucujici: "SA", Hodnota: 0x2, Datum: "13.11.2018", Mala: false},
			{Popis: "", Vyucujici: "SA", Hodnota: 0x2, Datum: "14.12.2018", Mala: false}},
			Final: 0x2},

		znamkyPredmet{Nazev: "Občanská\u00a0 nauka", Zkratka: "ON", Znamky: []znamka{
			{Popis: "diktát", Vyucujici: "SC", Hodnota: 0x2, Datum: "05.10.2018", Mala: false},
			{Popis: "test", Vyucujici: "SU", Hodnota: 0x0, Datum: "12.11.2018", Mala: false},
			{Popis: "zkoušení", Vyucujici: "SU", Hodnota: 0x1, Datum: "16.11.2018", Mala: false},
			{Popis: "zpráva", Vyucujici: "SU", Hodnota: 0x1, Datum: "10.12.2018", Mala: true},
			{Popis: "test", Vyucujici: "SU", Hodnota: 0x2, Datum: "14.12.2018", Mala: false}},
			Final: 0x2},

		znamkyPredmet{Nazev: "Matematika", Zkratka: "M", Znamky: []znamka{
			{Popis: "", Vyucujici: "HR", Hodnota: 0x1, Datum: "08.10.2018", Mala: false},
			{Popis: "", Vyucujici: "HR", Hodnota: 0x2, Datum: "18.11.2018", Mala: false},
			{Popis: "", Vyucujici: "HR", Hodnota: 0x3, Datum: "18.11.2018", Mala: false},
			{Popis: "", Vyucujici: "HR", Hodnota: 0x0, Datum: "28.11.2018", Mala: false},
			{Popis: "", Vyucujici: "HR", Hodnota: 0x5, Datum: "01.01.2019", Mala: false}},
			Final: 0x4},

		znamkyPredmet{Nazev: "Fyzika", Zkratka: "F", Znamky: []znamka{
			{Popis: "", Vyucujici: "RK", Hodnota: 0x0, Datum: "28.09.2018", Mala: false},
			{Popis: "", Vyucujici: "RK", Hodnota: 0x2, Datum: "12.10.2018", Mala: false},
			{Popis: "", Vyucujici: "RK", Hodnota: 0x1, Datum: "21.10.2018", Mala: false},
			{Popis: "", Vyucujici: "RK", Hodnota: 0x0, Datum: "03.12.2018", Mala: false},
			{Popis: "", Vyucujici: "RK", Hodnota: 0x2, Datum: "20.01.2019", Mala: false}},
			Final: 0x3},

		znamkyPredmet{Nazev: "Chemie", Zkratka: "Ch", Znamky: []znamka{
			{Popis: "test", Vyucujici: "VB", Hodnota: 0x5, Datum: "28.09.2018", Mala: false},
			{Popis: "", Vyucujici: "BO", Hodnota: 0x5, Datum: "15.01.2019", Mala: false}},
			Final: 0x5},

		znamkyPredmet{Nazev: "Informační a komunikační technologie", Zkratka: "IT", Znamky: []znamka{
			{Popis: "", Vyucujici: "MZ", Hodnota: 0x5, Datum: "02.10.2018", Mala: false},
			{Popis: "práce v hodině", Vyucujici: "MZ", Hodnota: 0x5, Datum: "09.10.2018", Mala: true},
			{Popis: "", Vyucujici: "MZ", Hodnota: 0x5, Datum: "23.10.2018", Mala: false},
			{Popis: "", Vyucujici: "SB", Hodnota: 0x5, Datum: "16.10.2018", Mala: false},
			{Popis: "", Vyucujici: "SB", Hodnota: 0x5, Datum: "09.11.2018", Mala: false}},
			Final: 0x5},

		znamkyPredmet{Nazev: "Tělesná\u00a0 výchova", Zkratka: "TV", Znamky: []znamka{
			{Popis: "", Vyucujici: "HB", Hodnota: 0x0, Datum: "19.11.2018", Mala: false}},
			Final: 0x0},

		znamkyPredmet{Nazev: "Počítačové aplikace", Zkratka: "PA", Znamky: []znamka{
			{Popis: "", Vyucujici: "SM", Hodnota: 0x1, Datum: "08.11.2018", Mala: false},
			{Popis: "", Vyucujici: "SM", Hodnota: 0x1, Datum: "12.12.2018", Mala: false},
			{Popis: "", Vyucujici: "SM", Hodnota: 0x5, Datum: "18.10.2018", Mala: false},
			{Popis: "", Vyucujici: "SM", Hodnota: 0x5, Datum: "15.11.2018", Mala: false}},
			Final: 0x2},

		znamkyPredmet{Nazev: "Základy elektrotechniky", Zkratka: "ZE", Znamky: []znamka{
			{Popis: "Lineární rovnice - písemka", Vyucujici: "BP", Hodnota: 0x0, Datum: "14.10.2018", Mala: false},
			{Popis: "elektrostatika - písemka", Vyucujici: "BP", Hodnota: 0x2, Datum: "23.10.2018", Mala: false},
			{Popis: "elektrická kapacita - zkoušení", Vyucujici: "BP", Hodnota: 0x2, Datum: "03.12.2018", Mala: false}},
			Final: 0x2},

		znamkyPredmet{Nazev: "Praxe", Zkratka: "DC", Znamky: []znamka{
			{Popis: "Cvičná práce - konstrukční montáž robotického auta LEGO", Vyucujici: "FI", Hodnota: 0x1, Datum: "11.10.2018", Mala: false},
			{Popis: "Řešení příkladů - programování kostky LEGO EV3", Vyucujici: "FI", Hodnota: 0x3, Datum: "18.10.2018", Mala: true},
			{Popis: "vičná práce - zápis a provedení funkčního SW řízení osv. LED, motoru", Vyucujici: "FI", Hodnota: 0x0, Datum: "17.11.2018", Mala: false},
			{Popis: "Opakování - závěrečná praktická práce", Vyucujici: "FI", Hodnota: 0x1, Datum: "06.12.2018", Mala: false},
			{Popis: "", Vyucujici: "SI", Hodnota: 0x0, Datum: "22.01.2019", Mala: false}},
			Final: 0x2},

		znamkyPredmet{Nazev: "Chování", Zkratka: "", Znamky: []znamka(nil), Final: 0x5},
	}
	testGetRequest(t, mockJecnak, "score/student",
		"_test/znamky.html-test", func() interface{} { return mockJecnak.GetZnamky() }, expected)
}

func TestJecnak_GetJecnitel(t *testing.T) {
	expected := Jecnitel{
		Jmeno:   "Josif Vissarionovič Stalin",
		Zkratka: "St", UzivatelskeJmeno: "stalin",
		Email:             "stalin@example.com",
		SoukromyEmail:     "stalin@stalingrad.example.com",
		Telefon:           "123 456 789 a linka 666",
		Kabinet:           "K66 - Stalingrad",
		TridniUcitel:      "E1a",
		KonzultacniHodiny: "kdykoliv po dohodě",
		ImgUrl:            "/img/thumbnail/IMG-stalinjosif.JPG",
	}
	testGetRequest(t, mockJecnak, "ucitel/ST",
		"_test/jecnitel-st.html-test", func() interface{} { return mockJecnak.GetJecnitel("ST") }, expected)
}

func TestJecnak_GetSeznamJecnitelu(t *testing.T) {
	expected := SeznamJecnitelu{
		SeznamJecniteluRow{Jmeno: "Adolf Hitler", Zkratka: "Hi"},
		SeznamJecniteluRow{Jmeno: "Josef Vissarionovič Stalin", Zkratka: "St"},
		SeznamJecniteluRow{Jmeno: "Benito Mussolini", Zkratka: "Ms"},
		SeznamJecniteluRow{Jmeno: "Eva Braunová", Zkratka: "Braunova"},
		SeznamJecniteluRow{Jmeno: "Mao Ce-tung", Zkratka: "Mo"},
		SeznamJecniteluRow{Jmeno: "Vladimir Iljič Lenin", Zkratka: "Le"},
		SeznamJecniteluRow{Jmeno: "Adolf Eichmann", Zkratka: "Ei"},
		SeznamJecniteluRow{Jmeno: "Heinrich Himmler", Zkratka: "Hm"},
	}
	testGetRequest(t, mockJecnak, "ucitel",
		"_test/seznamjecnitelu.html-test", func() interface{} { return mockJecnak.GetSeznamJecnitelu() }, expected)
}
