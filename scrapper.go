package spamtoputocorreos

import (
	"crypto/tls"
	"github.com/gocolly/colly"
	"net/http"
)

func FactoryCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (HTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.WithTransport(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}})
	return c
}
