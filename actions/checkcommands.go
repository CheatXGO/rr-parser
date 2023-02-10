package actions

import (
	"log"
	"strings"

	"github.com/CheatXGO/rr-parser/database"
	"github.com/CheatXGO/rr-parser/structures"
)

var Raids = map[string]string{"VoG": "Vault of Glass\n",
	"DSC":  "Deep Stone Crypt\n",
	"GoS":  "Garden of Salvation\n",
	"LW":   "Last Wish\n",
	"CoS":  "Crown of Sorrow\n",
	"SotP": "Scourge of the Past\n",
	"SoS":  "Leviathan, Spire of Stars\n",
	"EoW":  "Leviathan, Eater of Worlds\n",
	"Levi": "Leviathan\n",
	"VotD": "Vow of the Disciple\n",
	"KF":   "King's Fall\n"}

func Checker(bot structures.BotMessage, cfg structures.Config) {
	//messageID := bot.Message.MessageId    // bot message ID
	username := bot.Message.From.Username // chat user name
	userID := bot.Message.From.Id         // chat user id
	//isbot := bot.Message.From.IsBot       // user are bot? true false
	chatID := bot.Message.Chat.Id // chat id (if userID=chatID then it's private message)
	//mesDate := bot.Message.Date           // bot message date
	//userText := bot.Message.Text          // text from user
	usercommand := strings.Split(bot.Message.Text, " ")
	switch strings.ToLower(usercommand[0]) {
	case "/help":
		text := "rr - check player stats (format: /rr raid nickname)\n\n" +
			"reg - register your Destiny 2 nickname to your tg profile (format: /reg destiny2_nickname)\n\n" +
			"upd - update your registered Destiny 2 nickname (format: /upd destiny2_nickname)\n\n" +
			"my - check your stats (format: /my raid). Only for registered users.\n\n" +
			"lists - raids abbreviation\n\n" +
			"help - commands and raids abbreviation\n\n"
		go func() {
			err := Sender(text, chatID, cfg)
			if err != nil {
				e := structures.MyError{Func: "/help", Err: "can't send commands"}
				log.Fatalln(e.Error())
			}
		}()
	case "/start":
		text := "Hello, " + username + "! Type /help to commands and raids abbreviation\n\n"
		go func() {
			err := Sender(text, chatID, cfg)
			if err != nil {
				e := structures.MyError{Func: "/start", Err: "can't send commands"}
				log.Fatalln(e.Error())
			}
		}()
	case "/lists":
		text := ""
		for i, j := range Raids {
			text = text + i + " - " + j
		}
		go func() {
			err := Sender(text, chatID, cfg)
			if err != nil {
				e := structures.MyError{Func: "/lists", Err: "can't send commands"}
				log.Fatalln(e.Error())
			}
		}()
	case "/rr":
		if cap(usercommand) >= 3 && usercommand[1] != "" {
			go func() {
				err := Sender(CheckStats(usercommand[1], usercommand[2], true, cfg), chatID, cfg)
				if err != nil {
					e := structures.MyError{Func: "/rr CheckStats", Err: "can't resolve Bungie answer"}
					log.Fatalln(e.Error())
				}
			}()
		} else {
			go func() {
				err := Sender("Please fill right raid and D2 nickname after /rr", chatID, cfg)
				if err != nil {
					e := structures.MyError{Func: "/rr Sender", Err: "can't send commands"}
					log.Fatalln(e.Error())
				}
			}()
		}
	case "/my":
		if cap(usercommand) >= 2 && usercommand[1] != "" {
			go func() {
				res, bol := database.Dblookup(bot, cfg)
				err := Sender(CheckStats(usercommand[1], res, bol, cfg), chatID, cfg)
				if err != nil {
					e := structures.MyError{Func: "/my CheckStats", Err: "can't resolve Bungie answer"}
					log.Fatalln(e.Error())
				}
			}()
		} else {
			go func() {
				err := Sender("Please fill right raid abbreviation after /my", chatID, cfg)
				if err != nil {
					e := structures.MyError{Func: "/my Sender", Err: "can't send commands"}
					log.Fatalln(e.Error())
				}
			}()
		}
	case "/reg":
		if cap(usercommand) >= 2 && usercommand[1] != "" {
			go func() {
				err := Sender(database.Dbinsert(usercommand[1], bot, cfg), chatID, cfg)
				if err != nil {
					e := structures.MyError{Func: "/reg Sender", Err: "can't create DB record"}
					log.Fatalln(e.Error())
				}
			}()
		} else {
			go func() {
				err := Sender("Please fill Destiny 2 nickname after /reg", chatID, cfg)
				if err != nil {
					e := structures.MyError{Func: "/reg Sender", Err: "can't send commands"}
					log.Fatalln(e.Error())
				}
			}()
		}
	case "/upd":
		if cap(usercommand) >= 2 && usercommand[1] != "" {
			go func() {
				err := Sender(database.Dbupdate(usercommand[1], bot, cfg), chatID, cfg)
				if err != nil {
					e := structures.MyError{Func: "/upd Sender", Err: "can't update DB record"}
					log.Fatalln(e.Error())
				}
			}()
		} else {
			go func() {
				err := Sender("Please fill Destiny 2 nickname after /upd for change", chatID, cfg)
				if err != nil {
					e := structures.MyError{Func: "/upd Sender", Err: "can't send commands"}
					log.Fatalln(e.Error())
				}
			}()
		}
	default:
		if userID == chatID { //private message
			go func() {
				err := Sender("Check command syntax! /help", chatID, cfg)
				if err != nil {
					e := structures.MyError{Func: "json.Unmarshal &botText", Err: "can't decode botText"}
					log.Fatalln(e.Error())
				}
			}()
		}
	}
}
