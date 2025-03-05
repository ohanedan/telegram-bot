package telegram_bot

import (
	"context"
	"regexp"
	"telegram_bot/config"
	"telegram_bot/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/teilomillet/gollm"
)

type MessageEntry struct {
	User    string
	Message string
}

func printAiReplyMessages(l *logger.Logger, chatIDMap config.ChatIntMap) {
	for _, chat := range chatIDMap {
		for _, reply := range chat.AiReplies {
			log := l.Sprintf("Chat: green{%v} Prompt: green{%v}",
				chat.ChatName, reply.Prompt)
			if reply.If != "" {
				log = l.Sprintf("%v If: green{%v}", log, reply.If)
			}
			if reply.Regex != "" {
				log = l.Sprintf("%v Regex: green{%v}", log, reply.Regex)
			}
			if reply.ResponseRegex != "" {
				log = l.Sprintf("%v Response Regex: green{%v}", log, reply.ResponseRegex)
			}
			l.Println("AI REPLY", log)
		}
	}
}

func replyAiMessages(update tgbotapi.Update, l *logger.Logger,
	bot *tgbotapi.BotAPI, chatIDMap config.ChatIntMap, errChan chan error) {

	if update.Message == nil {
		l.Println("AI REPLY", "Incoming update is not a message")
		return
	}

	chat, ok := chatIDMap[update.Message.Chat.ID]
	if !ok {
		l.Println("AI REPLY", "No config file for green{%v}",
			update.Message.Chat.ID)
		return
	}

	for _, reply := range chat.AiReplies {
		replyAiMessage(reply, update, l, bot, chat, errChan)
	}
}

var chatHistory = make(map[int64][]gollm.PromptMessage)

func replyAiMessage(reply *config.AiReply, update tgbotapi.Update, l *logger.Logger,
	bot *tgbotapi.BotAPI, chat *config.Chat, errChan chan error) {

	ctx := context.Background()

	ok, err := regexp.MatchString(reply.Regex, update.Message.Text)
	if err != nil {
		errChan <- err
		return
	}

	if !ok {
		l.Println("AI REPLY", "Regex: green{%v} not match. Message: green{%v}",
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
			l.Println("AI REPLY", "Script: green{%v} not ok for incoming message.",
				reply.If)
			return
		}
	}

	llm, err := gollm.NewLLM(
		gollm.SetProvider("ollama"),
		gollm.SetAPIKey("none"),
		gollm.SetModel(reply.Model),
		gollm.SetLogLevel(gollm.LogLevelInfo),
	)
	if err != nil {
		errChan <- err
		return
	}

	chatHistory[chat.ChatID] = append(chatHistory[chat.ChatID], gollm.PromptMessage{
		Role:    update.Message.From.UserName,
		Content: update.Message.Text,
	})

	prompt := gollm.NewPrompt("answer last message according to system prompt", gollm.WithMessages(chatHistory[chat.ChatID]), gollm.WithSystemPrompt(reply.PromptText, "ephemeral"))

	response, err := llm.Generate(ctx, prompt)
	if err != nil {
		errChan <- err
		return
	}

	re, err := regexp.Compile(reply.ResponseRegex)
	if err != nil {
		errChan <- err
		return
	}

	matches := re.FindAllStringSubmatch(response, -1)
	if len(matches) != 1 {
		l.Println("AI REPLY", "Message: not sent. Chat: yellow{%v}. More than one matches. Response: yellow{%v}",
			chat.ChatName, response)
		return
	} else {
		message := matches[0][1]
		if message == "" {
			l.Println("AI REPLY", "Message empty. Chat: green{%v}", chat.ChatName)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)

		chatHistory[chat.ChatID] = append(chatHistory[chat.ChatID], gollm.PromptMessage{
			Role:    "assistant",
			Content: message,
		})

		l.Println("AI REPLY", "Message: green{%v} sent. Chat: green{%v}",
			message, chat.ChatName)
	}
}
