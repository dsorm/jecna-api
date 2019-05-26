package jecnaapi

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req"
	"net/http"
)

func (jecnak *Jecnak) getHttpResponse(url string) *http.Response {
	response, err := req.Get(url, req.Header{
		"User-Agent": "Mozilla/5.0",
		"Cookie":     jecnak.LoginSession + "; role=student;",
	})
	handleErr(err)
	return response.Response()
}

func (jecnak *Jecnak) getGoqueryDocument(url string) *goquery.Document {
	document, err := goquery.NewDocumentFromReader(jecnak.getHttpResponse(url).Body)
	handleErr(err)
	return document
}
