package spamtoputocorreos

import (
	"github.com/ervitis/spamtoputocorreos/models"
	"github.com/gocolly/colly"
	"log"
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

func GetTokens(c *colly.Collector, tokens *models.Tokens) {
	if tokens == nil {
		tokens = new(models.Tokens)
	}

	c.OnHTML("input[type=hidden]", func(element *colly.HTMLElement) {
		if element.Attr("name") == "_csrf" {
			tokens.Csrf = element.Attr("value")
		}
	})

	c.OnResponse(func(response *colly.Response) {
		cookies := c.Cookies(response.Request.URL.String())
		if len(cookies) == 1 {
			tokens.Session = cookies[0].Value
		}
	})

	err := c.Visit(urlHome)
	if err != nil {
		log.Println("error visiting home page", err)
	}
}
