// Hlavni soubor, neobsahujeinstall skoro nic mimo konstanty siteURL pro pripad
// zmeneni domeny Jecny a zasadni funkce isLogged()
package jecnaapi

import (
	"github.com/imroc/req"
	"golang.org/x/net/html"
	"io"
	"net/http"
)

var siteURL = "https://www.spsejecna.cz/" // s lomenem na konci

// Pokusi se o prihlaseni, zjisti, zda-li je platne, a ulozi session token do Jecnaka
// Je nutne zavolat pred vetsinou funkci ktera by na webu vyzadovala prihlaseni (napr. suplarchy)
// Pozor, session tokeny maji omezenou dobu platnosti, a taky ztraci platnost
// prihlasenim na portal z jineho zarizeni/programu
// Je nutne vyplnit pole LoginName a LoginPass v Jecnakovi pred volanim teto funkce
func (jecnak *Jecnak) IsLogged() bool {
	if len(jecnak.LoginSession) < 1 {
		// prihlasit uzivatele
		jecnak.LoginValid = false
		// pripojeni na stranky jecne pro ziskani session cookie
		resp, err := http.Get(siteURL) // imroc/req tu zpusobuje nejakou divnou chybu, proto pouzivam primo http
		handleErr(err)
		cookie := resp.Header.Get("Set-Cookie")
		jecnak.LoginSession = cookie
		// nastaveni post parametru a headeru
		loginParam := req.Param{
			"user": jecnak.LoginName,
			"pass": jecnak.LoginPass,
		}
		header := req.Header{
			"User-Agent": "Mozilla/5.0",
			"Cookie":     cookie,
		}
		// samotne prihlaseni
		response, err := req.Post(siteURL+"user/login", header, loginParam)
		// zjisteni chyb
		if err != nil {
			handleErr(err)
		} else {
			// tokenizerem zjistit, zdali je uzivatel prihlasen
			tokenizer := html.NewTokenizer(response.Response().Body)
			for {
				tokenType := tokenizer.Next()
				if tokenType == html.TextToken {
					if string(tokenizer.Text()) == "SPŠE Ječná - Novinky" {
						jecnak.LoginValid = true
						return true
					} else if string(tokenizer.Text()) == "SPŠE Ječná - Přihlášení se nezdařilo" {
						jecnak.LoginValid = false
						return true
					}
				} else if tokenizer.Err() == io.EOF {
					jecnak.LoginValid = false
					return true
				}
			}
		}

	} else {
		// overit, zda je stale prihlasen
		header := req.Header{
			"User-Agent": "Mozilla/5.0",
			"Cookie":     jecnak.LoginSession,
		}
		// samotne prihlaseni
		response, err := req.Post(siteURL+"user/login", header)
		// zjisteni chyb
		if err != nil {
			handleErr(err)
		} else {
			// tokenizerem zjistit, zdali je uzivatel prihlasen
			tokenizer := html.NewTokenizer(response.Response().Body)
			for {
				tokenType := tokenizer.Next()
				if tokenType == html.TextToken {
					if string(tokenizer.Text()) == "SPŠE Ječná - Novinky" {
						jecnak.LoginValid = true
						return true
					} else if string(tokenizer.Text()) == "SPŠE Ječná - Přihlášení se nezdařilo" {
						jecnak.LoginValid = false
						return true
					}
				} else if tokenizer.Err() == io.EOF {
					jecnak.LoginValid = false
					return true
				}
			}
		}
	}
	return jecnak.LoginValid
}
