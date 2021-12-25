package spamtoputocorreos

import (
	"crypto/tls"
	"github.com/gocolly/colly"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	GlobalSignalHandler = make(chan os.Signal, 1)

	DataToken *Tokens
)

func InitSignalHandler() {
	signal.Notify(GlobalSignalHandler, syscall.SIGTERM, os.Interrupt)
}

func NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout:   15 * time.Second,
	}
}

func getParams(text string) (paramsMap map[string]string) {
	match := reg.FindStringSubmatch(text)

	paramsMap = make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

func GetTokens(c *colly.Collector) *Tokens {
	token := new(Tokens)

	c.OnHTML("input[type=hidden]", func(element *colly.HTMLElement) {
		if element.Attr("name") == "_csrf" && token.Csrf == "" {
			token.Csrf = element.Attr("value")
		}
	})

	c.OnResponse(func(response *colly.Response) {
		cookies := c.Cookies(response.Request.URL.String())
		if len(cookies) == 1 {
			token.Session = cookies[0].Value
		}
	})

	_ = c.Visit(urlHome)

	return token
}
