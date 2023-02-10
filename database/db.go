package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/CheatXGO/rr-parser/structures"
	_ "github.com/lib/pq"
)

func Dbinsert(userd2name string, bot structures.BotMessage, cfg structures.Config) string {
	var text string
	var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Pip, cfg.Pport, cfg.Plg, cfg.Ppw, cfg.Pdb)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		e := structures.MyError{Func: "sql.Open", Err: "Failed to connect with postgres DB"}
		log.Fatalln(e.Error())
	}
	defer db.Close()
	rows, err := db.Query("select * from tgd2users WHERE idtelegr=$1 OR tgnick=$2", bot.Message.From.Id, bot.Message.From.Username)
	if err != nil {
		e := structures.MyError{Func: "db.Query", Err: "Failed to select from DB"}
		log.Fatalln(e.Error())
	}
	defer rows.Close()
	tuser := structures.TelegUser{}
	for rows.Next() {
		err = rows.Scan(&tuser.Id, &tuser.Idtelegr, &tuser.Tgnick, &tuser.D2nick)
		if err != nil {
			e := structures.MyError{Func: "rows.Scan", Err: "Failed to encode query"}
			log.Fatalln(e.Error())
		}
	}
	if tuser.Idtelegr == bot.Message.From.Id && tuser.Tgnick == bot.Message.From.Username {
		text = "You already have registered as " + tuser.D2nick + "\nUse /upd to change nickname"
		return text
	} else {
		_, err = db.Exec("insert into tgd2users (idtelegr, tgnick, d2nick) values ($1, $2, $3)", bot.Message.From.Id, bot.Message.From.Username, userd2name)
		if err != nil {
			e := structures.MyError{Func: "db.Exec", Err: "Failed to insert into DB"}
			log.Fatalln(e.Error())
		}
		text = "Succesfully registered as " + userd2name
	}
	return text
}

func Dbupdate(userd2name string, bot structures.BotMessage, cfg structures.Config) string {
	var text string
	var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Pip, cfg.Pport, cfg.Plg, cfg.Ppw, cfg.Pdb)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		e := structures.MyError{Func: "sql.Open", Err: "Failed to connect with postgres DB"}
		log.Fatalln(e.Error())
	}
	defer db.Close()
	rows, err := db.Query("select * from tgd2users WHERE idtelegr=$1 OR tgnick=$2", bot.Message.From.Id, bot.Message.From.Username)
	if err != nil {
		e := structures.MyError{Func: "db.Query", Err: "Failed to select from DB"}
		log.Fatalln(e.Error())
	}
	defer rows.Close()
	tuser := structures.TelegUser{}
	for rows.Next() {
		err = rows.Scan(&tuser.Id, &tuser.Idtelegr, &tuser.Tgnick, &tuser.D2nick)
		if err != nil {
			e := structures.MyError{Func: "rows.Scan", Err: "Failed to encode query"}
			log.Fatalln(e.Error())
		}
	}
	if tuser.Idtelegr == bot.Message.From.Id && tuser.Tgnick == bot.Message.From.Username {
		_, err = db.Exec("update tgd2users set d2nick=$1 where id=$2", userd2name, tuser.Id)
		if err != nil {
			e := structures.MyError{Func: "db.Exec", Err: "Failed to insert into DB"}
			log.Fatalln(e.Error())
		}
		text = "Username succesfully updated: " + userd2name
	} else {
		text = "You are not registered yet\nUse /reg  " + userd2name + "to register"
		return text
	}
	return text
}

func Dblookup(bot structures.BotMessage, cfg structures.Config) (string, bool) {
	var text string
	var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Pip, cfg.Pport, cfg.Plg, cfg.Ppw, cfg.Pdb)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		e := structures.MyError{Func: "sql.Open", Err: "Failed to connect with postgres DB"}
		log.Fatalln(e.Error())
	}
	defer db.Close()
	rows, err := db.Query("select * from tgd2users WHERE idtelegr=$1 OR tgnick=$2", bot.Message.From.Id, bot.Message.From.Username)
	if err != nil {
		e := structures.MyError{Func: "db.Query", Err: "Failed to select from DB"}
		log.Fatalln(e.Error())
	}
	defer rows.Close()
	tuser := structures.TelegUser{}
	for rows.Next() {
		err = rows.Scan(&tuser.Id, &tuser.Idtelegr, &tuser.Tgnick, &tuser.D2nick)
		if err != nil {
			e := structures.MyError{Func: "rows.Scan", Err: "Failed to encode query"}
			log.Fatalln(e.Error())
		}
	}
	if tuser.Idtelegr == bot.Message.From.Id && tuser.Tgnick == bot.Message.From.Username {
		text = tuser.D2nick
	} else {
		text = "You are not registered yet\nUse /reg to register"
		return text, false
	}
	return text, true
}
