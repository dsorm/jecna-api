package main

import (
	"fmt"
	jecnaapi "github.com/david-sorm/jecna-api"
)

func main() {
	// Vytvoříme proměnou a vyplníme v ní jméno a heslo
	jecnak := jecnaapi.Jecnak{
		LoginName: "sorm",     // nahradit vlastním přihlašovacím jménem
		LoginPass: "02656161", // nahradit vlastním heslem
	}

	// Touto funkcí se přihlasíme a zároveň ověříme platnost údajů
	jePrihlasen := jecnak.IsLogged()

	// Pokud je výstup z funkce nepravdivý, znamená to, že údaje nejsou platné
	if !jePrihlasen {
		fmt.Println("Neplatné přihlašovací údaje!")
		return
	}

	// Stahneme udaje
	jecnitel := jecnak.GetJecnitel("bp")

	// A zobrazime je
	jecnitel.Print()

	// Muzeme treba pokracovat realistickym vyuzitim Ječná API - děláním velmi vtipných vtipů
	fmt.Printf("\n\nStudent A: Kdo je tvůj nejoblíbenější učitel?\n")
	fmt.Printf("Student B: Samozřejmě že %v!\n", jecnitel.Jmeno)
}
