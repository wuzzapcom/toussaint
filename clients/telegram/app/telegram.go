package app

import (
	"log"
	"toussaint/clients/telegram/app/cache"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewTelegram(token string, debug bool) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = debug

	log.Printf("telegram: authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	tg := &Telegram{
		bot:          bot,
		updateConfig: u,
		shouldStop:   make(chan bool, 0),
		stopped:      make(chan bool, 0),
	}

	return tg, nil
}

type Telegram struct {
	bot          *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig

	shouldStop chan bool
	stopped    chan bool
}

func (tg *Telegram) Start() error {

	notifier := NewNotifier(tg)
	notifier.Start()
	defer notifier.Stop()
	log.Printf("started notifier")

	updateChan, err := tg.bot.GetUpdatesChan(tg.updateConfig)
	if err != nil {
		return err
	}

	for {
		select {
		case upd := <-updateChan:
			go func() {
				answer, err := cache.HandleMessage(upd.Message)
				if err != nil {
					log.Printf("[ERR] cache.HandleMessage: %+v", err)
				}

				err = tg.answer(upd.Message.Chat.ID, answer)
				if err != nil {
					log.Printf("[ERR] tg.answer: %+v", err)
				}
			}()

		case <-tg.shouldStop:
			tg.stopped <- true
			return nil
		}
	}
}

func (tg *Telegram) answer(recipient int64, message string) error {
	msg := tgbotapi.NewMessage(recipient, message)

	_, err := tg.bot.Send(msg)
	return err
}

func (tg *Telegram) Stop() {
	tg.shouldStop <- true
	<-tg.stopped
}
