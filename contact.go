package spamtoputocorreos

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	urlContact = "https://www.adtpostales.com/webauth/adtPostales/public/processContact"
)

func NewContactData() *ContactData {
	return &ContactData{
		Name:               "Victor Martin",
		Phone:              "685990843",
		Email:              "vitomarti@gmail.com",
		Category:           "ENVIOS",
		InquiryDescription: "NO HE RECIBIDO LA NOTIFICACIÓN DE LLEGADA DE MI ENVÍO",
		Inquiry:            "NO_LLEGADA_ENVIO",
		RefCode:            "CN029952816JP",
		Query:              "Buenos días, ¿En qué fecha aproximada recibiré el aviso y el paquete? Un saludo",
		AcceptPrivacy:      true,
	}
}

func Contact(client *http.Client, tokens *Tokens, contactData *ContactData) error {
	form := map[string]string{
		"_csrf":                tokens.Csrf,
		"name":                 contactData.Name,
		"phone":                contactData.Phone,
		"email":                contactData.Email,
		"categoria":            contactData.Category,
		"categoriaDescripcion": "",
		"solicitudDescripcion": contactData.InquiryDescription,
		"solicitud":            contactData.Inquiry,
		"nEnvio":               contactData.RefCode,
		"consulta":             contactData.Query,
		"acceptPrivacy":        strconv.FormatBool(contactData.AcceptPrivacy),
		"tokenReCaptcha":       tokens.Captcha,
	}

	headers := map[string]string{
		"User-Agent":                "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (HTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language":           "es-ES,es;q=0.8,en-US;q=0.5,en;q=0.3",
		"Accept-Encoding":           "gzip, deflate, br",
		"Origin":                    "https://www.adtpostales.com",
		"Connection":                "keep-alive",
		"Cookie":                    "UISESSION=" + tokens.Session,
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
			return err
		}
	}

	_, err := w.CreateFormFile("fichero", "")
	if err != nil {
		return err
	}

	_, err = w.CreateFormFile("fichero2", "")
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, urlContact, body)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	req.Header.Add("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error doing request, code %d status %s", resp.StatusCode, resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanLines)
	var buffer bytes.Buffer
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return err
	}

	doc.Find("#errorDiv")
	return nil
}
