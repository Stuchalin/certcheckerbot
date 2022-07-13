package main

import (
	"certcheckerbot/botprocessing"
	"log"
	"os"
)

func main() {
	myBot, err := botprocessing.NewBot(os.Getenv("BOT_KEY"))
	if err != nil {
		log.Panic(err)
	}

	errorsBot := botprocessing.StartProcessing(myBot)

	for {
		select {
		case err := <-errorsBot:
			log.Printf("Bot error message: %s", err)
		}
	}

}
