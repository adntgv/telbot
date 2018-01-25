package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/net/context"

	speech "cloud.google.com/go/speech/apiv1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
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
		cmd := exec.Command("ffmpeg", "-i", "mes.ogg", "mes.wav")
		err = cmd.Run()
		if err != nil {
			log.Printf("%s", err)
		}

		rec, err := recognize("mes.wav")
		if err != nil {
			log.Printf("%s", err)
		}

		message := ""
		for _, result := range rec.Results {
			for _, alt := range result.Alternatives {
				fmt.Printf("\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
				message = alt.Transcript
			}
		}

		err = os.Remove("mes.ogg")
		if err != nil {
			log.Printf("%s", err)
		}

		err = os.Remove("mes.wav")
		if err != nil {
			log.Printf("%s", err)
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}

func recognize(file string) (*speechpb.RecognizeResponse, error) {
	ctx := context.Background()

	// [START init]
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// [END init]

	// [START request]
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	resp, err := client.Recognize(ctx, &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:     speechpb.RecognitionConfig_LINEAR16,
			LanguageCode: "ru-RU",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: data},
		},
	})
	// [END request]
	return resp, err
}
