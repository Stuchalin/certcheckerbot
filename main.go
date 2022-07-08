package main

import (
	"crypto/tls"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_KEY"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			command := update.Message.Text
			command = strings.Trim(command, " ")
			if command[:1] == "/" {
				msgText, err := commandProcessing(command)

				if err != nil {
					log.Println("Error in Dial", err)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID

		}
	}
}

// Processing known commands
func commandProcessing(command string) (string, error) {
	// Parse command and attributes
	i := strings.Index(command, " ")
	cmd := command
	attr := ""
	if i != -1 {
		cmd = command[:i]
		attr = command[i+1:]
	}
	// Remove double spaces
	for ok := true; ok; ok = strings.Contains(attr, "  ") {
		attr = strings.Replace(attr, "  ", " ", -1)
	}
	// Execute commands
	switch cmd {
	case "/help":
		return "Here will be print a help message", nil
	case "/check":
		if attr == "" {
			return "You must specify the URL. Format: \n\t /check [www.checkURL1.com www.checkURL2.com ...]. Use space to check few URLs.", nil
		}
		return getCertsInfo(attr, false), nil
	default:
		return "Use /help command", nil
	}
}

func getCertsInfo(URLs string, printFullChain bool) string {
	UrlArr := strings.Split(URLs, " ")
	result := ""
	for _, url := range UrlArr {
		result += getCertInfo(url, printFullChain)
	}
	return result
}

func getCertInfo(URL string, printFullChain bool) string {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", URL+":443", conf)
	if err != nil {
		log.Println("Error in Dial", err)
		return fmt.Sprintf("Cannot check cert from URL %s. Error: %e\n\n", URL, err)
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	result := ""
	for _, cert := range certs {
		if !printFullChain && cert.IsCA {
			continue
		}
		result += fmt.Sprintf("DNSNames: %s\n", cert.DNSNames)
		result += fmt.Sprintf("Issuer Name: %s\n", cert.Issuer)
		result += fmt.Sprintf("Expiry: %s \n", cert.NotAfter.Format("2006-02-02"))
		result += fmt.Sprintf("Common Name: %s \n", cert.Issuer.CommonName)

	}
	return result + "\n"
}
