package main

import (
	"net/http"
	"log"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("516390842:AAH7Sd4t_0J5gxyEYHUCbU9D9jwVt7gZdd4")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Voice == nil {
			continue
		} 
		var f tgbotapi.FileConfig
		f.FileID = update.Message.Voice.FileID
		file, err := bot.GetFile(f)
		if err != nil{
			log.Printf("%s",err)
		}
		resp, err := http.Get(file.FilePath)
		if err != nil {
			log.Printf("%s",err)
		}
		defer resp.Body.Close()
		
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
