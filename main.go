package main

import (
	"certcheckerbot/botprocessing"
	"certcheckerbot/storage/sqlite3"
	"log"
	"os"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	db, err := sqlite3.NewController(dbPath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Dispose()

	myBot, err := botprocessing.NewBot(os.Getenv("BOT_KEY"), db)
	if err != nil {
		log.Panic(err)
	}

	errorsBot := myBot.StartProcessing()

	for {
		select {
		case err := <-errorsBot:
			log.Printf("Bot error message: %s", err)
		}
	}

}
