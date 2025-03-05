package telegram_bot

import (
	"telegram_bot/config"
	"telegram_bot/logger"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func RunBot(config *config.Scheme) error {

	l := &logger.Logger{
		Disabled: !config.EnableLogs,
	}

	bot, err := tgbotapi.NewBotAPI(config.APIKey)
	if err != nil {
		return err
	}
	l.Println("RUN", "Authorized on account green{%s}", bot.Self.UserName)

	u := tgbotapi.NewUpdate(-1)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return nil
	}

	printReplyMessages(l, config.IDChats)

	errChan := make(chan error)

	go scheduledMessages(l, bot, config.Chats, errChan)

	lastActivityTime := time.Now()
	for {
		select {
		case update := <-updates:
			go replyMessages(update, l, bot, config.IDChats, errChan)
			lastActivityTime = time.Now()
		case err := <-errChan:
			return err
		default:
			if time.Now().Sub(lastActivityTime).Seconds() < 300 {
				continue
			}
			l.Println("RUN", "No activity for last 5 minutes.")
			lastActivityTime = time.Now()
		}
		time.Sleep(100 * time.Millisecond)
	}
}
