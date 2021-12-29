package main

import (
	"fmt"
	"github.com/ervitis/spamtoputocorreos"
	"github.com/ervitis/spamtoputocorreos/repo"
	"log"
)

func init() {
	repo.LoadDBConfig()
	spamtoputocorreos.LoadTelegramConfig()
	spamtoputocorreos.LoadCustomsData()
}

func main() {
	crawler := spamtoputocorreos.FactoryCollector()

	spamtoputocorreos.DataToken = spamtoputocorreos.GetTokens(crawler)
	spamtoputocorreos.DataToken.Captcha = spamtoputocorreos.CustomsData.Captcha

	db := repo.New(&repo.DBConfig)

	svc := spamtoputocorreos.NewCustomsTracerService(db)
	st, err := svc.GetStatus(spamtoputocorreos.DataToken, spamtoputocorreos.CustomsData.RefCode)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(st)
}
