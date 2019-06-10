// Tento soubor obsahuje funkce pro vypis vsech datovych typu,
// ktere vyflusnou funkce ze souboru get.go
// Jsou urceny spise pro debugging a demonstraci nez
// pouziti v "produkcnich" prostredich, holt lepsi neco nez nic
package jecnaapi

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var stdout io.Writer = os.Stdout

// Vypise do konzole obsah datove struktury SuplarchObsah
// Pokud je do argumentu vlozen Jecnak, automaticky
// vypise pouze suplovani jeho tridy
// Pokud je do nej vlozen typ int (pozor, ne napr. int64 nebo uint8)
// nastavuje se tim uroven podrobnosti vypisu
// Mozne hodnoty intu: 1, 2
// Pokud je vlozena hodnota vyssi ci nizsi, automaticky se vybere ta nejblizsi
func (suplarchObsah SuplarchObsah) Print(a interface{}) {
	var details uint8
	switch a.(type) {
	case Jecnak:
		if len(a.(Jecnak).SkupinyTrida) >= 2 {
			tridaNalezena := false
			trida := strings.Split(a.(Jecnak).SkupinyTrida, ", ")[0]
			for _, v := range suplarchObsah {
				if strings.ToLower(v.trida) == strings.ToLower(trida) {
					tridaNalezena = true
					fmt.Fprintf(stdout, "Suplování třídy %v:", v.trida)
					for k2, v2 := range v.supl {
						fmt.Fprintf(stdout, "\n [%v]: %v", k2+1, MakeSuplovaniTextReadable(v2))
					}
					fmt.Fprintf(stdout, "\n")
				}
			}
			if !tridaNalezena {
				fmt.Println("Třída nebyla nalezena, vypisuji všechny...")
				details = 1
			} else {
				return
			}
		} else {
			if len(a.(Jecnak).SkupinyTrida) < 1 {
				fmt.Println("Ječňák nemá vyplněnou třídu (co třeba Jecnak.GetJecnakInfo()?), vypisuji všechny...")
			}
			details = 1
		}
	case int:
		if a.(int) < 1 {
			details = 1
		} else if a.(int) > 2 {
			details = 2
		} else {
			details = uint8(a.(int))
		}
	default:
		details = 1
	}

	switch details {
	case 1:
		for _, v := range suplarchObsah {
			fmt.Fprintf(stdout, "\nTřída: %v\nSuplování:", v.trida)
			for k2, v2 := range v.supl {
				if len(v2) > 0 {
					fmt.Fprintf(stdout, " [%v]: %v", k2, MakeSuplovaniTextReadable(v2))
				}
			}
			fmt.Fprintf(stdout, "\n")
		}
	case 2:
		for _, v := range suplarchObsah {
			fmt.Fprintf(stdout, "\nTřída: %v\nSuplování:", v.trida)
			for k2, v2 := range v.supl {
				fmt.Fprintf(stdout, " [%v]: %v\n", k2, MakeSuplovaniTextReadable(v2))
			}
		}
	}

}

// Vypise do konzole obsah datove struktury SuplarchSeznam
func (suplarchSeznam SuplarchSeznam) Print() {
	for k, v := range suplarchSeznam {
		fmt.Fprintf(stdout, "[%v] %v (%v)\n", k, v.name, v.url)
	}
}

// Vypise do konzole obsah datove struktury Prichody
// Pokud chcete vypsat i dny, kdy nebyl zaznamenan zadny prichod,
// nastavte argument "printEmptyRows" na true
func (prichody Prichody) Print(printEmptyRows bool) {
	for _, v := range prichody {
		if printEmptyRows && len(v.PrichodOdchod) < 1 {
			continue
		}
		fmt.Println(v.Datum+":", v.PrichodOdchod)
	}
}

// Vypise do konzole obsah datove struktury OmluvnyList
func (omluvnyList OmluvnyList) Print() {
	for _, v := range omluvnyList {
		fmt.Println(v.Datum+":", v.Text)
	}
}

// Vypise do konzole obsah datove struktury Jecnak
// Pokud chcete zobrazit i heslo k prihlaseni, nastavte
// argument "showPassword" na true
func (jecnak *Jecnak) Print(showPassword bool) {
	fmt.Fprintf(stdout, "-- Přihlašovací údaje --\n")
	fmt.Fprintf(stdout, "Přihlašovací jméno: %v\n", jecnak.LoginName)
	if showPassword {
		fmt.Fprintf(stdout, "Heslo: %v\n", jecnak.LoginPass)
	}
	fmt.Fprintf(stdout, "Session: %v\n", jecnak.LoginSession)
	fmt.Fprintf(stdout, "Je login platný?: ")
	if jecnak.LoginValid {
		fmt.Fprintf(stdout, "Ano\n")
	} else {
		fmt.Fprintf(stdout, "Ne\n")
	}
	if len(jecnak.CeleJmeno) < 1 {
		fmt.Fprintf(stdout, "Osobní údaje Ječňáka nenalezeny\n")
	} else {
		fmt.Fprintf(stdout, "-- Osobní údaje --\n")
		fmt.Fprintf(stdout, "Celé jméno: %v\n", jecnak.CeleJmeno)
		fmt.Fprintf(stdout, "Věk: %v\n", jecnak.Vek)
		fmt.Fprintf(stdout, "Narozeniny: %v\n", jecnak.Narozeniny)
		fmt.Fprintf(stdout, "Místo narození: %v\n", jecnak.MistoNarozeni)
		fmt.Fprintf(stdout, "Telefon: %v\n", jecnak.Telefon)
		fmt.Fprintf(stdout, "Trvalá adresa: \n")
		fmt.Fprintf(stdout, "  %v\n  %v %v\n", jecnak.TrvalaAdresaUlice, jecnak.TrvalaAdresaMesto, jecnak.TrvalaAdresaPSC)
		fmt.Fprintf(stdout, "Třída, skupiny: %v\n", jecnak.SkupinyTrida)
		fmt.Fprintf(stdout, "Číslo v třídním výkazu: %v\n", jecnak.CisloVTridnimVykazu)
		fmt.Fprintf(stdout, "Soukromý email: %v\n", jecnak.SoukromyEmail)
		fmt.Fprintf(stdout, "Školní email: %v\n", jecnak.SkolniEmail)
	}

}

// Vypise do konzole obsah datove struktury Rozvrh
// Je to trosku vizualne rozbity, protoze je mi to nejak jedno,
// ale obsah se vypisuje spravne
func (rozvrh *Rozvrh) Print() {
	space := "           "
	for _, v := range rozvrh.Casy {
		fmt.Fprintf(stdout, v.CisloHodiny)
		fmt.Fprintf(stdout, space)
	}

	fmt.Fprintf(stdout, "\n")

	for _, v := range rozvrh.Casy {
		fmt.Fprintf(stdout, v.Doba)
		fmt.Fprintf(stdout, space)
	}

	fmt.Fprintf(stdout, "\n")

	for i := 1; i <= 5; i++ {
		var dnes []RozvrhHodina
		switch i {
		case 1:
			dnes = rozvrh.Pondeli
		case 2:
			dnes = rozvrh.Utery
		case 3:
			dnes = rozvrh.Streda
		case 4:
			dnes = rozvrh.Ctvrtek
		case 5:
			dnes = rozvrh.Patek
		}
		for _, v := range dnes {
			fmt.Fprintf(stdout, v.ZkratkaHodiny)
			fmt.Fprintf(stdout, space)

		}
		fmt.Fprintf(stdout, "\n")
	}
}

// Vypise do konzole obsah datove struktury Znamky
// Do argumentu "details" lze zadat uroven podrobnosti vypisu
// Mozne hdonoty: 1, 2, 3
// V pripade zadani jine hodnoty bude pouzita ta nejblizsi mozna k zadane
func (znamky Znamky) Print(details uint8) {
	if details > 3 {
		details = 3
	} else if details == 0 {
		details = 1
	}
	for _, v := range znamky {
		fmt.Println(v.Nazev + ": ")
		for _, v2 := range v.Znamky {
			if v2.Hodnota != 0 {
				switch details {
				case 1:
					fmt.Fprintf(stdout, " %d", v2.Hodnota)
				case 2:
					datum, _ := time.Parse("02.01.2006", v2.Datum)
					if datum.Year() == 0001 {
						fmt.Fprintf(stdout, " [%d] %v\n", v2.Hodnota, v2.Popis)
					} else {
						fmt.Fprintf(stdout, " [%d] %v (%v)\n", v2.Hodnota, v2.Popis, datum.Format("2.1.2006"))
					}
				case 3:
					datum, _ := time.Parse("02.01.2006", v2.Datum)
					if datum.Year() == 0001 {
						fmt.Fprintf(stdout, " [%d] %v: %v\n", v2.Hodnota, v2.Vyucujici, v2.Popis)
					} else {
						fmt.Fprintf(stdout, " [%d] %v: %v (%v)\n", v2.Hodnota, v2.Vyucujici, v2.Popis, datum.Format("2.1.2006"))
					}
				}
			} else {
				switch details {
				case 1:
					fmt.Fprintf(stdout, " N")
				case 2:
					fmt.Fprintf(stdout, " [N] %v (%v)\n", v2.Popis, v2.Datum)
				case 3:
					fmt.Fprintf(stdout, "[N] %v: %v (%v)\n", v2.Vyucujici, v2.Popis, v2.Datum)
				}
			}
		}
		if v.Final != 100 {
			if v.Final != 0 {
				switch details {
				case 1:
					fmt.Fprintf(stdout, " [%v]", v.Final)
				case 2, 3:
					fmt.Fprintf(stdout, " konečná známka: [%v]", v.Final)
				}
			} else {
				switch details {
				case 1:
					fmt.Fprintf(stdout, " [N]")
				case 2, 3:
					fmt.Fprintf(stdout, " konečná známka: [N]")
				}
			}

		}
		fmt.Fprintf(stdout, "\n\n")
	}
}

// Vypise do konzole obsah datove struktury Jecnitel
func (jecnitel *Jecnitel) Print() {
	fmt.Fprintf(stdout, "Jméno: %v (%v)\n", jecnitel.Jmeno, jecnitel.Zkratka)
	fmt.Fprintf(stdout, "Uživatelské jméno: %v\n", jecnitel.UzivatelskeJmeno)
	fmt.Fprintf(stdout, "Email: %v\n", jecnitel.Email)
	fmt.Fprintf(stdout, "Soukromý email: %v\n", jecnitel.SoukromyEmail)
	fmt.Fprintf(stdout, "Telefon: %v\n", jecnitel.Telefon)
	fmt.Fprintf(stdout, "Kabinet: %v\n", jecnitel.Kabinet)
	fmt.Fprintf(stdout, "Třídnictví: %v\n", jecnitel.TridniUcitel)
	fmt.Fprintf(stdout, "Konzultační hodiny: %v\n", jecnitel.KonzultacniHodiny)
	if len(jecnitel.ImgUrl) > 1 {
		fmt.Fprintf(stdout, "Fotka: %v\n", strings.TrimSuffix(siteURL, "/")+jecnitel.ImgUrl)
	}
}

// Vypise do konzole obsah datove struktury SeznamJecnitelu
func (seznamJecnitelu SeznamJecnitelu) Print() {
	for _, v := range seznamJecnitelu {
		fmt.Fprintf(stdout, "[%v] %v\n", v.Zkratka, v.Jmeno)
	}
}
