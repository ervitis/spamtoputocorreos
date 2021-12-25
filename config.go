package spamtoputocorreos

import "github.com/kelseyhightower/envconfig"

type (
	TelegramParams struct {
		Token  string `envconfig:"TELEGRAM_TOKEN"`
		UserID int64  `envconfig:"TELEGRAM_USER_ID"`
	}

	CustomsDataTrace struct {
		RefCode string `envconfig:"REF_CODE"`
		Captcha string `envconfig:"CAPTCHA_CODE"`
	}
)

var (
	TelegramConfig TelegramParams
	CustomsData    CustomsDataTrace
)

func LoadTelegramConfig() {
	envconfig.MustProcess("", &TelegramConfig)
}

func LoadCustomsData() {
	envconfig.MustProcess("", &CustomsData)
}
