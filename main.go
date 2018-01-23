package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	token := "516390842:AAH7Sd4t_0J5gxyEYHUCbU9D9jwVt7gZdd4"
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

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
		if err != nil {
			log.Printf("%s", err)
		}
		resp, err := http.Get("https://api.telegram.org/file/bot" + token + "/" + file.FilePath)
		if err != nil {
			log.Printf("%s", file.FilePath)
			log.Printf("%s", err)
		}
		defer resp.Body.Close()
		ogg, err := os.Create("mes.ogg")
		if err != nil {
			log.Printf("%s", err)
		}

		_, err = io.Copy(ogg, resp.Body)
		if err != nil {
			log.Printf("%s", err)
		}
		cmd := exec.Command("soundconverter", "mes.ogg", "-b", "-s", ".wav", "-m", "wav")
		err = cmd.Run()
		if err != nil {
			log.Printf("%s", err)
		}
		
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
