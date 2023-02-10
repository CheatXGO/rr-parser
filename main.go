package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/CheatXGO/rr-parser/actions"
	"github.com/CheatXGO/rr-parser/certs"
	"github.com/CheatXGO/rr-parser/structures"
	"github.com/CheatXGO/rr-parser/tgconnection"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg := structures.Config{}
	file, err := os.ReadFile("config.json")
	if err != nil {
		e := structures.MyError{Func: "io.ReadFile", Err: "can't read config.json"}
		log.Fatalln(e.Error())
	}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		e := structures.MyError{Func: "json.Unmarshal", Err: "problem with decoding config.json"}
		log.Fatalln(e.Error())
	}
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		e := structures.MyError{Func: "tgbotapi.NewBotAPI", Err: "can't initialize bot"}
		log.Fatalln(e.Error())
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	err = tgconnection.SetWebhook(bot, cfg.TelegramBotToken, cfg.TgBotIP, cfg.TgBotPort, certs.CERT)
	if err != nil {
		e := structures.MyError{Func: "tgconnection.SetWebhook", Err: "can't set webhook for bot"}
		log.Fatalln(e.Error())
	}
	message := func(w http.ResponseWriter, r *http.Request) {
		text, err := io.ReadAll(r.Body)
		if err != nil {
			e := structures.MyError{Func: "io.ReadAll", Err: "can't read bot message"}
			log.Fatalln(e.Error())
		}
		var botText structures.BotMessage
		err = json.Unmarshal(text, &botText)
		if err != nil {
			e := structures.MyError{Func: "main json.Unmarshal &botText", Err: "can't decode botText"}
			log.Fatalln(e.Error())
		}
		go actions.Checker(botText, cfg)
	}
	http.HandleFunc("/", message)
	err = http.ListenAndServeTLS(":8443", certs.CERT, certs.KEY, nil)
	if err != nil {
		e := structures.MyError{Func: "http.ListenAndServeTLS", Err: "http server don't start"}
		log.Fatalln(e.Error())
	}
}
