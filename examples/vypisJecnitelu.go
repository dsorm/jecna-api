package main

import "github.com/david-sorm/jecna-api"

func main() {
	// Nejdřív je potřeba udělat proměnou s typem Jecnak
	jecnak := jecnaapi.Jecnak{}

	// Teď můžeme zavolat funkci pro získání dat
	seznamJecnitelu := jecnak.GetSeznamJecnitelu()

	// Poté pro demonstraci tento list vypíšeme
	// Tato funkce vypisuje seznam ve tvaru [zkratka] jméno
	seznamJecnitelu.Print()
}
