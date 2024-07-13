package main

import (
	"os"
	"time"

	db "tg_bot/db"
	tg "tg_bot/telegrambot"
)

func main() {
	time.Sleep(1 * time.Minute)

	if os.Getenv("CREATE_TABLE") == "yes" {
		if os.Getenv("DB_SWITCH") == "on" {
			if err := db.CreateTables(); err != nil {
				panic(err)
			}
		}
	}

	time.Sleep(1 * time.Minute)

	// Запускаем бота
	tg.TelegramBot()
}
