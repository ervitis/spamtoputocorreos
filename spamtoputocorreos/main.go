package main

import (
	"fmt"
	"github.com/ervitis/spamtoputocorreos"
	"github.com/ervitis/spamtoputocorreos/bots"
	"github.com/ervitis/spamtoputocorreos/repo"
	"log"
	"os"
	"time"
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

	svc := spamtoputocorreos.NewCustomsTracerService(db)

	tb, err := bots.NewTelegramBot(&spamtoputocorreos.TelegramConfig, svc, db)
	if err != nil {
		log.Panicln(err)
	}

	ticker := time.NewTicker(4 * time.Hour)

	go func() {
		for {
			select {
			case <-spamtoputocorreos.GlobalSignalHandler:
				return
			case t := <-ticker.C:
				hasUpdate, err := svc.SearchTracerUpdatesAndUpdatesDB()
				if err != nil {
					log.Panicln(err)
				}
				if hasUpdate {
					_ = tb.SendNotification(fmt.Sprintf("Tick at %s. There is a new update! type the command `/latest` to see latest information", t))
				}
			}
		}
	}()

	go func(tb *bots.TelegramBot) {
		log.Panicln(tb.StartServer())
	}(tb)

	server := spamtoputocorreos.NewServer(port)
	server.Start()

}
