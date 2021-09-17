package actions

import (
	"bytes"
	"encoding/json"
	"log"
	"main/structures"
	"net/http"
)

func Sender(text string, chatID int, cfg structures.Config) error {
	if chatID != 0 { //check for user-edited messages
		reqBody := structures.SendMessage{Id: chatID, Text: text}
		reqBytes, err := json.Marshal(reqBody)
		if err != nil {
			e := structures.MyError{Fun: "json.Marshal &reqBody", Err: "can't encode answer for user"}
			log.Fatalln(e.Error())
		}
		res, err := http.Post("https://api.telegram.org/bot"+cfg.TelegramBotToken+"/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
		if err != nil {
			e := structures.MyError{Fun: "http.Post", Err: "can't send answer"}
			log.Fatalln(e.Error())
		}
		if res.StatusCode != http.StatusOK {
			e := structures.MyError{Fun: "http.StatusOK", Err: "answer send but not OK"}
			println(e.Error(), chatID, text)
		}
	}
	return nil
}
