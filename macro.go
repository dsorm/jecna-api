package jecnaapi

func (jecnak Jecnak) GetSuplarchAndPrint(n uint8) {
	jecnak.IsLogged()
	jecnak.GetJecnakInfo()
	SuplarchContent := jecnak.GetSuplarchContent(jecnak.GetSuplarchSeznam()[n])
	SuplarchContent.Print(jecnak)
}
