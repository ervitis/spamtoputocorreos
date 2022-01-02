package spamtoputocorreos

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ervitis/spamtoputocorreos/models"
	"github.com/ervitis/spamtoputocorreos/regtools"
	"github.com/ervitis/spamtoputocorreos/repo"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"regexp"
	"sort"
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
		db       repo.IRepository
	}

	ICustomsStatus interface {
		GetStatus(*models.Tokens, string) (*models.StatusTrace, error)
		SearchTracerUpdatesAndUpdatesDB() (bool, error)
	}
)

func NewCustomsTracerService(db repo.IRepository) ICustomsStatus {
	return &CustomsStatusTrace{
		scrapper: FactoryCollector(),
		db:       db,
	}
}

func (c *CustomsStatusTrace) GetStatus(tokens *models.Tokens, refCode string) (*models.StatusTrace, error) {
	form := make(map[string]string)
	form["envio.numEnvio"] = refCode
	form["tokenReCaptcha"] = tokens.Captcha
	form["_csrf"] = tokens.Csrf

	status := &models.StatusTrace{RefCode: refCode, Statuses: make([]*models.StatusData, 0)}

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
				data := regtools.GetParams(reg, text)
				st := strings.Split(data["Status"], splitSpaces)
				t, err := time.Parse(dateFormatTraceLayout, data["Date"])
				if err != nil {
					log.Fatal(err)
				}

				trace := &models.StatusData{
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

func (c *CustomsStatusTrace) SearchTracerUpdatesAndUpdatesDB() (bool, error) {
	statuses, err := c.GetStatus(DataToken, CustomsData.RefCode)
	if err != nil {
		return false, fmt.Errorf("error getting status from customs web: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	data, err := c.db.Get(ctx, statuses.RefCode)
	if err != nil {
		return false, fmt.Errorf("error getting data from db: %w", err)
	}

	lenStatuses := len(statuses.Statuses)
	lenData := len(data.Statuses)

	if lenData > 0 {
		sort.Slice(data.Statuses, func(i, j int) bool {
			return data.Statuses[i].Date.After(data.Statuses[j].Date)
		})
	}

	if lenStatuses > 0 {
		sort.Slice(statuses.Statuses, func(i, j int) bool {
			return statuses.Statuses[i].Date.After(statuses.Statuses[j].Date)
		})
	}

	indexData := 0
	if lenData > 0 {
		indexData = lenData - 1
	}

	indexStatuses := 0
	if lenStatuses > 0 {
		indexStatuses = lenStatuses - 1
	}

	log.Printf("indexes statuses %d, data %d. Len statuses %d, data %d", indexStatuses, indexData, lenStatuses, lenData)

	for _, v := range statuses.Statuses {
		log.Printf("statuses date %s", v.Date.Format(time.RFC822))
	}

	for _, v := range data.Statuses {
		log.Printf("data db date %s", v.Date.Format(time.RFC822))
	}

	if indexData > 0 && indexStatuses > 0 && data.Statuses[indexData].Date.After(statuses.Statuses[indexStatuses].Date) {
		log.Printf("No updates in the package %s :(", statuses.RefCode)
		return false, nil
	}

	log.Printf("latest date from data %s, latest date from scrap %s", data.Statuses[indexData].Date.Format(time.RFC3339), statuses.Statuses[indexStatuses].Date.Format(time.RFC3339))
	log.Printf("after? %v", data.Statuses[indexData].Date.After(statuses.Statuses[indexStatuses].Date)
	log.Printf("There is a new update for the package\n")

	if err := c.db.Save(ctx, statuses); err != nil {
		return false, fmt.Errorf("error saving data: %w", err)
	}
	return true, nil
}
