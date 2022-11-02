package main

import (
	"certcheckerbot/botprocessing"
	"certcheckerbot/scheduler"
	"certcheckerbot/storage"
	"certcheckerbot/storage/sqlite3"
	"encoding/json"
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

	var days []int
	envDays := os.Getenv("EXPIRY_DAYS")
	err2 := json.Unmarshal([]byte(envDays), &days)
	if err2 != nil {
		log.Printf("\nFail to convert array with notification days - %v. Set default notification days.\n", envDays)
		days = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 60, 90}
	}

	myBot, err := botprocessing.NewBot(os.Getenv("BOT_KEY"), db, debug)
	if err != nil {
		log.Panic(err)
	}

	usersDomainsChan := make(chan *storage.User, 100)
	errorsBot := myBot.StartProcessing(usersDomainsChan, days)

	go scheduler.InitScheduler(db, usersDomainsChan)

	for {
		select {
		case err := <-errorsBot:
			log.Printf("Bot error message: %s", err)
		}
	}

}
