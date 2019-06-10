package jecnaapi

// Stahne n-ty suplarch a zobrazi ho
// Demonstrace vyuziti Jecna API
// Vyzaduje platne vyplnene prihlasovaci udaje v Jecnakovi, a argument "n"
// Hodnota n = 0 znaci nejnovejsi suplarch, hodnota 1 predposledni,
// a tak dal az do 4 (limit seznamu suplarchu v Jecna API)
func (jecnak Jecnak) GetSuplarchAndPrint(n uint8) {
	// ziskame session token
	jecnak.IsLogged()
	// ziskame tridu Jecnaka, abychom nemuseli inteligentne vypisovat cely suplarch
	jecnak.GetJecnakInfo()
	// stahne seznam suplarchu, a z toho vytahne n-ty, ze ktereho potom stahne obsah a ulozi do
	// promenne SuplarchContent
	SuplarchContent := jecnak.GetSuplarchContent(jecnak.GetSuplarchSeznam()[n])
	// vypise do konzole
	SuplarchContent.Print(jecnak)
}
