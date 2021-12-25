package main

import (
	"github.com/ervitis/spamtoputocorreos"
	"github.com/ervitis/spamtoputocorreos/bots"
	"github.com/ervitis/spamtoputocorreos/repo"
	"log"
	"os"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	svc := spamtoputocorreos.NewCustomsTracerService()

	tb, err := bots.NewTelegramBot(&spamtoputocorreos.TelegramConfig, svc, db)
	if err != nil {
		log.Panicln(err)
	}

	go func(tb *bots.TelegramBot) {
		log.Panicln(tb.StartServer())
	}(tb)

	server := spamtoputocorreos.NewServer(port)
	server.Start()

}
