package spamtoputocorreos

import (
	"github.com/ervitis/spamtoputocorreos/models"
	"github.com/gocolly/colly"
	"os"
	"os/signal"
	"syscall"
)

var (
	GlobalSignalHandler = make(chan os.Signal, 1)

	DataToken *models.Tokens
)

func InitSignalHandler() {
	signal.Notify(GlobalSignalHandler, syscall.SIGTERM, os.Interrupt)
}

func GetTokens(c *colly.Collector) *models.Tokens {
	token := new(models.Tokens)

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
