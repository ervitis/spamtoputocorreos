package bots

import (
	"fmt"
	"github.com/ervitis/spamtoputocorreos"
	"github.com/ervitis/spamtoputocorreos/models"
	"github.com/ervitis/spamtoputocorreos/repo"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"sort"
	"strings"
	"time"
)

type (
	// TelegramBot represents a Telegram bot.
	TelegramBot struct {
		user           *tb.User
		bot            *tb.Bot
		traceService   spamtoputocorreos.ICustomsStatus
		contactService *spamtoputocorreos.ContactService
		db             repo.IRepository
	}
)

func NewTelegramBot(cfg *spamtoputocorreos.TelegramParams, contactService *spamtoputocorreos.ContactService, traceService spamtoputocorreos.ICustomsStatus, db repo.IRepository) (*TelegramBot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot:            b,
		traceService:   traceService,
		contactService: contactService,
		db:             db,
		user:           &tb.User{ID: cfg.UserID},
	}, nil
}

func (t *TelegramBot) SendNotification(msg string) error {
	_, err := t.bot.Send(t.user, msg)
	return err
}

func (t *TelegramBot) StartServer() error {
	go func() {
		select {
		case <-spamtoputocorreos.GlobalSignalHandler:
			t.bot.Stop()
		}
	}()

	t.bot.Handle("/amialive", t.handleHealthCheck)
	t.bot.Handle("/latest", t.handleGetLatestStatus)
	t.bot.Handle("/all", t.handleGetAllStatus)
	t.bot.Handle("/search", t.handleSearchUpdatesAndNotify)
	t.bot.Handle("/help", t.handleHelp)
	t.bot.Handle("/submit", t.handleSendQueryContact)

	t.bot.Start()
	return nil
}

func (t *TelegramBot) registerError(err error, msgDetail string) {
	log.Printf("Detected an error %s: %s\n", err, msgDetail)
	_, err = t.bot.Send(t.user, fmt.Sprintf("error: %s: %s", err, msgDetail))
	if err != nil {
		log.Printf("I could not send the error message to telegram: %s\n", err)
	}
}

func (t *TelegramBot) handleHelp(_ *tb.Message) {
	var dt models.InquiryDescriptionType
	var ct models.InquiryCategoryType

	helpMessage := fmt.Sprintf(`
HELP COMMANDS:

- /latest
- /amialive
- /all
- /search
- /submit <message>. The <message> field has to follow this format:
	<inquiry_category>-<inquiry_description>-<content>
	<inquiry_category> can be %s (TBD)
	<inquiry_description> can be
		ENVIOS: %s
	<content> inquiry content of message
`, strings.Join(ct.All(), ","), strings.Join(dt.All(), ","))
	_, _ = t.bot.Send(t.user, helpMessage)
}

func (t *TelegramBot) handleGetLatestStatus(_ *tb.Message) {
	statuses, err := t.traceService.GetStatus(spamtoputocorreos.DataToken, spamtoputocorreos.CustomsData.RefCode)
	if err != nil {
		t.registerError(err, "handleGetLatestStatus.GetStatus")
		return
	}

	if len(statuses.Statuses) == 0 {
		t.registerError(fmt.Errorf("something happened while scrapping web"), "check code in scrapping")
		return
	}

	msg := fmt.Sprintf("Latest status of the package %s:\n", statuses.RefCode)

	sort.Slice(statuses.Statuses, func(i, j int) bool {
		return statuses.Statuses[i].Date.After(statuses.Statuses[j].Date)
	})

	s := statuses.Statuses[0]

	msg += fmt.Sprintf("\t%s: %s - %s\n", s.Date.Format(time.RFC3339), s.Status, s.Detail)
	_, err = t.bot.Send(t.user, msg)
	if err != nil {
		t.registerError(err, "handleGetLatestStatus send data of package")
	}
}

func (t *TelegramBot) handleGetAllStatus(_ *tb.Message) {
	statuses, err := t.traceService.GetStatus(spamtoputocorreos.DataToken, spamtoputocorreos.CustomsData.RefCode)
	if err != nil {
		t.registerError(err, "handleGetLatestStatus.GetStatus")
		return
	}

	if len(statuses.Statuses) == 0 {
		t.registerError(fmt.Errorf("something happened while scrapping web"), "check code in scrapping")
		return
	}

	msg := fmt.Sprintf("Statuses of the package %s:\n", statuses.RefCode)
	for _, v := range statuses.Statuses {
		msg += fmt.Sprintf("\t%s: %s - %s\n", v.Date.Format(time.RFC3339), v.Status, v.Detail)
	}
	_, err = t.bot.Send(t.user, msg)
	if err != nil {
		t.registerError(err, "handleGetAllStatus send data of package")
	}
}

func (t *TelegramBot) handleSearchUpdatesAndNotify(_ *tb.Message) {
	// check in database and notify if there is a new update
	hasUpdate, err := t.traceService.SearchTracerUpdatesAndUpdatesDB()
	if err != nil {
		t.registerError(err, "handleSearchUpdatesAndNotify")
		return
	}

	if !hasUpdate {
		return
	}
	t.handleGetLatestStatus(nil)
}

func (t *TelegramBot) handleHealthCheck(_ *tb.Message) {
	_, err := t.bot.Send(t.user, "I am alive!")
	if err != nil {
		log.Printf("error sending message in healthcheck %s\n", err)
	}
}

func (t *TelegramBot) handleSendQueryContact(m *tb.Message) {
	inquiryBody := &models.InquiryBodyData{}
	var err error
	inquiryBody, err = inquiryBody.Marshal(m.Payload)
	if err != nil {
		log.Println("Error handling contact to customs in message", err)
		_, _ = t.bot.Send(t.user, "An error happened sending the query to customs: %s", err)
		return
	}

	contactData := spamtoputocorreos.NewContactData(inquiryBody)

	if err := t.contactService.Contact(spamtoputocorreos.DataToken, contactData); err != nil {
		log.Println("Error in contact service", err)
		_, _ = t.bot.Send(t.user, fmt.Sprintf("An error happened sending the query to customs: %s", err))
		return
	}

	_, _ = t.bot.Send(t.user, fmt.Sprintf("Message to customs sent. Data sent: %s", inquiryBody.String()))
}
