package main

import (
	"certcheckerbot/botprocessing"
	"certcheckerbot/storage/sqlite3"
	"log"
	"os"
	"strconv"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	db, err := sqlite3.NewController(dbPath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Dispose()

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Printf("\nConvert debug value error - %v. Default debug set to false.\n", err)
		debug = false
	}
	myBot, err := botprocessing.NewBot(os.Getenv("BOT_KEY"), db, debug)
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
