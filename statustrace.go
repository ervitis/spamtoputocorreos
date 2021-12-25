package spamtoputocorreos

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	urlTrace = "https://www.adtpostales.com/webauth/adtPostales/public/seguimientoEnvio"
	urlHome  = "https://www.adtpostales.com/webauth/adtPostales/"

	dateFormatTraceLayout = "2/01/2006 15:04:05"

	regexPattern = `(?P<Date>\d{1,2}\/\d{1,2}\/\d{4}\s\d{1,2}:\d{2}:\d{2})\s+(?P<Status>.*)`

	splitSpaces = `       `
)

var (
	reg = regexp.MustCompile(regexPattern)
)

type (
	CustomsStatusTrace struct {
		scrapper *colly.Collector
	}
)

func NewCustomsTracerService() *CustomsStatusTrace {
	return &CustomsStatusTrace{
		scrapper: FactoryCollector(),
	}
}

func (c *CustomsStatusTrace) GetStatus(tokens *Tokens, refCode string) (*StatusTrace, error) {
	form := make(map[string]string)
	form["envio.numEnvio"] = refCode
	form["tokenReCaptcha"] = tokens.Captcha
	form["_csrf"] = tokens.Csrf

	status := &StatusTrace{RefCode: refCode, Statuses: make([]*StatusData, 0)}

	c.scrapper.OnRequest(func(request *colly.Request) {
		err := c.scrapper.SetCookies(urlTrace, []*http.Cookie{
			{
				Name:  "UISESSION",
				Value: tokens.Session,
			},
			{
				Name:  "msg_cookie_ADT",
				Value: strconv.FormatBool(false),
			},
			{
				Name:  "cookie_ADT_tecnica",
				Value: strconv.FormatBool(false),
			},
		})
		if err != nil {
			log.Println(err)
		}
	})

	c.scrapper.OnResponse(func(response *colly.Response) {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))
		if err != nil {
			panic(err)
		}

		doc.Find(`#listadoTraza tbody`).Each(func(i int, selection *goquery.Selection) {
			selection.Find("tr").Each(func(j int, line *goquery.Selection) {
				text := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line.Text(), "\n", ""), "\t", ""))
				data := getParams(text)
				st := strings.Split(data["Status"], splitSpaces)
				t, err := time.Parse(dateFormatTraceLayout, data["Date"])
				if err != nil {
					log.Fatal(err)
				}

				trace := &StatusData{
					Date:   t,
					Status: strings.TrimSpace(st[0]),
					Detail: strings.TrimSpace(st[1]),
				}

				status.Statuses = append(status.Statuses, trace)
			})
		})
	})

	err := c.scrapper.Post(urlTrace, form)
	if err != nil {
		return nil, err
	}

	return status, nil
}
