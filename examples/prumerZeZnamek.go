package main

import (
	"fmt"
	"github.com/david-sorm/jecna-api"
)

func main() {
	// Vytvoříme proměnou a vyplníme v ní jméno a heslo
	jecnak := jecnaapi.Jecnak{
		LoginName: "sorm", // nahradit vlastním přihlašovacím jménem
		LoginPass: "0",    // nahradit vlastním heslem
	}

	// Touto funkcí se přihlasíme a zároveň ověříme platnost údajů
	jePrihlasen := jecnak.IsLogged()

	// Pokud je výstup z funkce nepravdivý, znamená to, že údaje nejsou platné
	if !jePrihlasen {
		fmt.Println("Neplatné přihlašovací údaje!")
		return
	}

	// Stáhneme známky
	znamky := jecnak.GetZnamky()

	// Typ jecnaapi.Znamky je synonymem pro jecnaapi.znamkyPredmet[], tudíž se jedná pouze o slice předmětů
	// Proto můžeme na něj použít range
	var soucet uint8
	var pocetZnamek uint8
	for _, predmet := range znamky { // pro každý předmět
		pocetZnamek, soucet = 0, 0
		for _, znamka := range predmet.Znamky { // pro kazdou znamku z daneho predmetu
			// N-ka mají hodnotu 0, těma se nezabývame
			if znamka.Hodnota == 0 {
				continue
			}
			soucet += znamka.Hodnota
			pocetZnamek++
		}
		// Nulou se moc dělit nedá
		if pocetZnamek == 0 {
			continue
		}
		fmt.Printf("[%v]: %v (znamek: %v, soucet hodnot: %v)\n",
			predmet.Nazev,
			// soucet a pocetZnamek převedeme na float pro floatový výsledek
			float32(soucet)/float32(pocetZnamek),
			pocetZnamek,
			soucet,
		)
	}
}
