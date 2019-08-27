# Ječná API
[![Build Status](https://travis-ci.org/david-sorm/jecna-api.svg?branch=master)](https://travis-ci.org/david-sorm/jecna-api)

Neoficiální knihovna pro získání dat z portálu spsejecna.cz

Slovník:
 - ječnitel = učitel
 - ječňák = přihlášený uživatel
 
Funkce:
 - kontrola přihlašovacích údajů
 - seznam ječnitelů (i bez přihlášení)
 - informace o učiteli (bez přihlášení nejsou úplné)
 - známky
 - rozvrh 
 - informace o ječňákovi
 - omluvný list
 - příchody a odchody
 - seznam posledních pěti suplarchů
 - parsing dat ze suplarchu
 - jednoduchý výpis všech datových struktur pro debugging

Známé chyby:
 - rozvrh neumí pracovat s hodinami, které mají více skupin žáků
 - u známek nefunguje rozdělení hodin (např. praxe a teorie), jsou smíchány k sobě
 - I přes veškerou snahu zblbuvzdornit parsing suplarchů a cest k nim, nemusí tyto funkce fungovat, ale v tomto případě 
  za to může externí faktor - neschopnost dodržovat jednotnou cestu (např. značení školního roku) a obsah suplarchů 
  (kompletní změna uspořádání náležitostí, nemyslím změnu tříd atd.) od jistého pana zástupce. Kdyby aspoň nezapomínal 
  dávat suplarchy do novinek, tak bych nemusel dělat takovéto šaškárny.
  
Příklad použití:
```go
// Program pro výpis seznamu ječnitelů bez přihlašovacích údajů
package main

import "github.com/david-sorm/jecna-api"

func main() {
	// Nejdřív je potřeba udělat proměnou s typem Jecnak
	jecnak := jecnaapi.Jecnak{}

	// Teď můžeme zavolat funkci pro získání dat
	seznamJecnitelu := jecnak.GetSeznamJecnitelu();

	// Poté pro demonstraci tento list vypíšeme
	// Tato funkce vypisuje seznam ve tvaru [zkratka] jméno
	seznamJecnitelu.Print();
}
```
Tento a ještě další okomentované přiklady jsou obsaženy ve složce `examples`