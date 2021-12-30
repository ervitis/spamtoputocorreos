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

	CustomsDataContact struct {
		Telephone string `envconfig:"DATA_TELEPHONE"`
		Name      string `envconfig:"DATA_NAME"`
		Email     string `envconfig:"DATA_EMAIL"`
	}
)

var (
	TelegramConfig TelegramParams
	CustomsData    CustomsDataTrace
	ContactData    CustomsDataContact
)

func LoadContactData() {
	envconfig.MustProcess("", &ContactData)
}

func LoadTelegramConfig() {
	envconfig.MustProcess("", &TelegramConfig)
}

func LoadCustomsData() {
	envconfig.MustProcess("", &CustomsData)
}
