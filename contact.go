package spamtoputocorreos

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ervitis/spamtoputocorreos/models"
	"github.com/gocolly/colly"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	urlContact = "https://www.adtpostales.com/webauth/adtPostales/public/processContact"
)

type (
	ContactService struct {
		scrapper *colly.Collector
	}
)

func NewContactService() *ContactService {
	return &ContactService{scrapper: FactoryCollector()}
}

func NewContactData(inquiryData *models.InquiryBodyData) *models.ContactData {
	return &models.ContactData{
		Name:               ContactData.Name,
		Phone:              ContactData.Telephone,
		Email:              ContactData.Email,
		Category:           inquiryData.Category,    // InquiryCategoryType
		InquiryDescription: inquiryData.Description, // "NO HE RECIBIDO LA NOTIFICACIÓN DE LLEGADA DE MI ENVÍO",
		InquiryCategory:    inquiryData.Type,        // InquiryDescriptionType
		RefCode:            CustomsData.RefCode,
		Query:              inquiryData.Query, // "Buenos días, ¿En qué fecha aproximada recibiré el aviso y el paquete? Un saludo",
		AcceptPrivacy:      true,
	}
}

func (cs *ContactService) Contact(tokens *models.Tokens, contactData *models.ContactData) error {
	form := map[string]string{
		"_csrf":                tokens.Csrf,
		"name":                 contactData.Name,
		"phone":                contactData.Phone,
		"email":                contactData.Email,
		"categoria":            contactData.Category.Value(),
		"categoriaDescripcion": "",
		"solicitudDescripcion": contactData.InquiryDescription,
		"solicitud":            contactData.InquiryCategory.Value(),
		"nEnvio":               contactData.RefCode,
		"consulta":             contactData.Query,
		"acceptPrivacy":        strconv.FormatBool(contactData.AcceptPrivacy),
		"tokenReCaptcha":       tokens.Captcha,
	}

	headers := map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":           "es-ES,es;q=0.8,en-US;q=0.5,en;q=0.3",
		"Accept-Encoding":           "gzip, deflate, br",
		"Origin":                    "https://www.adtpostales.com",
		"Connection":                "keep-alive",
		"Upgrade-Insecure-Requests": "1",
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "same-origin",
		"Sec-Fetch-User":            "?1",
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	for key, val := range form {
		if err := w.WriteField(key, val); err != nil {
			return fmt.Errorf("contact service: setting cookies: %w", err)
		}
	}

	_, err := w.CreateFormFile("fichero", "")
	if err != nil {
		return fmt.Errorf("contact service: setting fichero: %w", err)
	}

	_, err = w.CreateFormFile("fichero2", "")
	if err != nil {
		return fmt.Errorf("contact service: setting fichero2: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("contact service: closing body build: %w", err)
	}

	cs.scrapper.OnRequest(func(request *colly.Request) {
		for k, v := range headers {
			request.Headers.Add(k, v)
		}

		if err := cs.scrapper.SetCookies(urlContact, []*http.Cookie{
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
		}); err != nil {
			log.Println(err)
		}

		request.Headers.Set("Content-Type", w.FormDataContentType())
	})

	var errResp error
	cs.scrapper.OnResponse(func(response *colly.Response) {
		query, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))
		if err != nil {
			log.Println("querying from body response error: ", err)
			return
		}
		errMsg := query.Find(`div[id="errorDiv"] > ul > li`).Text()
		if errMsg != "" {
			log.Println("response from web page when contacting customs: ", errMsg)
			errResp = fmt.Errorf("response from page: %s", errMsg)
		}
	})

	if err := cs.scrapper.PostRaw(urlContact, body.Bytes()); err != nil {
		errResp = fmt.Errorf("error doing request to contact page: %w", err)
	}

	return errResp
}
