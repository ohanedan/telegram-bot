package telegram_bot

import (
	"regexp"
	"strings"
	"telegram_bot/config"
	"telegram_bot/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func printReplyMessages(l *logger.Logger, chatIDMap config.ChatIntMap) {
	for _, chat := range chatIDMap {
		for _, reply := range chat.Replies {
			log := l.Sprintf("Chat: green{%v} Message: green{%v}",
				chat.ChatName, reply.Message)
			if reply.If != "" {
				log = l.Sprintf("%v If: green{%v}", log, reply.If)
			}
			if reply.Regex != "" {
				log = l.Sprintf("%v Regex: green{%v}", log, reply.Regex)
			}
			l.Println("REPLY", log)
		}
	}
}

func replyMessages(update tgbotapi.Update, l *logger.Logger,
	bot *tgbotapi.BotAPI, chatIDMap config.ChatIntMap, errChan chan error) {

	if update.Message == nil {
		l.Println("REPLY", "Incoming update is not a message")
		return
	}

	chat, ok := chatIDMap[update.Message.Chat.ID]
	if !ok {
		l.Println("REPLY", "No config file for green{%v}",
			update.Message.Chat.ID)
		return
	}

	for _, reply := range chat.Replies {
		replyMessage(reply, update, l, bot, chat, errChan)
	}
}

func replyMessage(reply *config.Reply, update tgbotapi.Update, l *logger.Logger,
	bot *tgbotapi.BotAPI, chat *config.Chat, errChan chan error) {

	ok, err := regexp.MatchString(reply.Regex, update.Message.Text)

	if err != nil {
		errChan <- err
		return
	}

	if !ok {
		l.Println("REPLY", "Regex: green{%v} not match. Message: green{%v}",
			reply.Regex, update.Message.Text)
		return
	}

	if reply.If != "" {
		ok, err := runScript(update.Message, reply.If)
		if err != nil {
			errChan <- err
			return
		}

		if !ok {
			l.Println("REPLY", "Script: green{%v} not ok for incoming message.",
				reply.If)
			return
		}
	}

	msgToSend := reply.Message
	msgToSend = strings.ReplaceAll(msgToSend, "(name)", update.Message.From.FirstName)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgToSend)
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)

	l.Println("REPLY", "Message: green{%v} sent. Chat: green{%v}",
		msgToSend, chat.ChatName)
}
