package models

import "time"

type (
	ContactData struct {
		Name               string
		Phone              string
		Email              string
		Category           string
		InquiryDescription string
		Inquiry            string
		RefCode            string
		Query              string
		AcceptPrivacy      bool
	}

	StatusTrace struct {
		RefCode  string
		Statuses []*StatusData
	}

	StatusData struct {
		Date   time.Time
		Status string
		Detail string
	}

	Tokens struct {
		Captcha string
		Csrf    string
		Session string
	}

	ReqStatusTraceBody struct {
		Code           string `json:"envio.numEnvio"`
		TokenRecaptcha string `json:"tokenRecaptcha"`
	}
)
