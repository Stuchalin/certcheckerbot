package main

import (
	"certcheckerbot/botprocessing"
	"certcheckerbot/storage/sqlite3"
	"log"
	"os"
	"strconv"
	"time"
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

	go initScheduler()

	for {
		select {
		case err := <-errorsBot:
			log.Printf("Bot error message: %s", err)
		}
	}

}

//Initialize the scheduler for every hour on the border of the next hour
func initScheduler() {
	now := time.Now()
	duration := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location()).Sub(now)
	log.Printf("Init scheduler after %v", duration)
	select {
	case <-time.After(duration):
		log.Println("Scheduler initialised!")
		for tick := range time.Tick(time.Hour) {
			log.Println(tick)
		}
	}
}
