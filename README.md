# speech2text_ru
Telegram bot, that recieves audio message and responds with text message with transcribed words.
Uses Google Cloud Speech API for recognition and ffmpeg for file conversion.
For it to work you have to:
+ have Google Speech API registered.
+ get token for your bot.
+ FFMPEG installed

Fetching:
$ go get github.com/right-hearted/telbot
$ cd {PATH_TO}/telbot
$ export GOOGLE_APPLICATION_CREDENTIALS="Path to credentials"
$ go run main.go
