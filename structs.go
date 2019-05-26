package jecnaapi

type SuplarchTrida struct {
	trida string
	supl  []string
}

type SuplarchObsah []SuplarchTrida

type SuplarchSeznam []Suplarch

type Suplarch struct {
	name string
	url  string
}

type Prichody []prichodyRow

type prichodyRow struct {
	Datum         string
	PrichodOdchod string
}

type OmluvnyList []omluvnyListRow

type omluvnyListRow struct {
	Datum string
	Text  string
}

type Rozvrh struct {
	Casy    []RozvrhCas
	Pondeli []RozvrhHodina
	Utery   []RozvrhHodina
	Streda  []RozvrhHodina
	Ctvrtek []RozvrhHodina
	Patek   []RozvrhHodina
}

type RozvrhHodina struct {
	UcitelZkratka string
	ZkratkaHodiny string
	Ucebna        string
	Trida         string
	Skupina       string
}

type RozvrhCas struct {
	CisloHodiny string
	Doba        string
}

type Znamky []znamkyPredmet

type znamkyPredmet struct {
	Nazev   string
	Zkratka string
	Znamky  []znamka
	Final   uint8
}

type znamka struct {
	Popis     string
	Vyucujici string
	Hodnota   uint8
	Datum     string
	Mala      bool
}

type SeznamJecnitelu []SeznamJecniteluRow

type SeznamJecniteluRow struct {
	Jmeno   string
	Zkratka string
}

type Jecnitel struct {
	// ucitel
	Jmeno             string
	Zkratka           string
	UzivatelskeJmeno  string
	Email             string
	SoukromyEmail     string
	Telefon           string
	Kabinet           string
	TridniUcitel      string
	KonzultacniHodiny string
	ImgUrl            string
	// Rozvrh             *Rozvrh // planovana funkce
}

type Jecnak struct {
	// zak
	// Login - login_name a login_pass musi byt vyplneny, pote IsLogged() pro ziskani sessionu nebo overeni loginu
	LoginName    string
	LoginPass    string
	LoginSession string
	LoginValid   bool
	// Osobni udaje - GetJecnakInfo()
	CeleJmeno           string
	Vek                 uint8
	Narozeniny          string
	MistoNarozeni       string
	Telefon             uint32
	TrvalaAdresaUlice   string
	TrvalaAdresaMesto   string
	TrvalaAdresaPSC     string
	SkupinyTrida        string
	CisloVTridnimVykazu uint8
	SoukromyEmail       string
	SkolniEmail         string
}
