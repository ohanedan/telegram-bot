package telegram_bot

import (
	"errors"
	"fmt"
	"telegram_bot/config"
	"telegram_bot/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
)

func scheduledMessages(l *logger.Logger, bot *tgbotapi.BotAPI,
	chatMap config.ChatStringMap, errChan chan error) {

	for _, chat := range chatMap {
		for _, message := range chat.ScheduledMessages {
			registerScheduledMessage(message, chat, l, bot, errChan)
		}
	}

	<-gocron.Start()
}

func registerScheduledMessage(mes *config.ScheduledMessage, chat *config.Chat,
	l *logger.Logger, bot *tgbotapi.BotAPI, errChan chan error) {

	if len(mes.Days) == 0 {
		scheduleForEveryDay(mes, chat, l, bot, errChan)
		return
	}

	for _, day := range mes.Days {
		scheduleForSpesificDay(day, mes, chat, l, bot, errChan)
	}
}

func scheduleForSpesificDay(day string, mes *config.ScheduledMessage, chat *config.Chat,
	l *logger.Logger, bot *tgbotapi.BotAPI, errChan chan error) {

	for _, when := range mes.When {
		switch day {
		case "Monday":
			gocron.Every(1).Monday().At(when).Do(sendScheduledMessage,
				l, bot, chat.ChatName, chat.ChatID, mes.Message, errChan)
		case "Tuesday":
			gocron.Every(1).Tuesday().At(when).Do(sendScheduledMessage,
				l, bot, chat.ChatName, chat.ChatID, mes.Message, errChan)
		case "Wednesday":
			gocron.Every(1).Wednesday().At(when).Do(sendScheduledMessage,
				l, bot, chat.ChatName, chat.ChatID, mes.Message, errChan)
		case "Thursday":
			gocron.Every(1).Thursday().At(when).Do(sendScheduledMessage,
				l, bot, chat.ChatName, chat.ChatID, mes.Message, errChan)
		case "Friday":
			gocron.Every(1).Friday().At(when).Do(sendScheduledMessage,
				l, bot, chat.ChatName, chat.ChatID, mes.Message, errChan)
		case "Saturday":
			gocron.Every(1).Saturday().At(when).Do(sendScheduledMessage,
				l, bot, chat.ChatName, chat.ChatID, mes.Message, errChan)
		case "Sunday":
			gocron.Every(1).Sunday().At(when).Do(sendScheduledMessage,
				l, bot, chat.ChatName, chat.ChatID, mes.Message, errChan)
		default:
			errChan <- errors.New(fmt.Sprintf("%v is not a valid day.", day))
		}

		l.Println("SCHEDULE", "Message: green{%v} Time: green{%v} "+
			"Chat: green{%v} Day: green{%v}",
			mes.Message, when, chat.ChatName, day)
	}
}

func scheduleForEveryDay(mes *config.ScheduledMessage, chat *config.Chat,
	l *logger.Logger, bot *tgbotapi.BotAPI, errChan chan error) {
	for _, when := range mes.When {
		gocron.Every(1).Day().At(when).Do(sendScheduledMessage,
			l, bot, chat.ChatName, chat.ChatID, mes.Message, errChan)

		l.Println("SCHEDULE", "Message: green{%v} Time: green{%v} Chat: green{%v}",
			mes.Message, when, chat.ChatName)
	}
}

func sendScheduledMessage(l *logger.Logger, bot *tgbotapi.BotAPI,
	chat_name string, chat_id int64, msg string, errChan chan error) {

	message := tgbotapi.NewMessage(chat_id, msg)
	_, err := bot.Send(message)
	if err != nil {
		errChan <- err
		return
	}

	l.Println("SCHEDULE", "Message: green{%v} sent. Chat: green{%v}",
		msg, chat_name)
}
