package botprocessing

import (
	"certcheckerbot/certinfo"
	"certcheckerbot/storage"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

type Bot struct {
	BotAPI *tgbotapi.BotAPI
	db     storage.UsersConfig
}

func NewBot(botKey string, db storage.UsersConfig, debug bool) (*Bot, error) {
	botApi, err := tgbotapi.NewBotAPI(botKey)
	if err != nil {
		return nil, err
	}

	botApi.Debug = debug

	log.Printf("Authorized on account %s", botApi.Self.UserName)

	bot := Bot{
		BotAPI: botApi,
		db:     db,
	}

	return &bot, nil
}

func (bot *Bot) StartProcessing() chan error {

	errorsChan := make(chan error, 10)

	go bot.startProcessing(errorsChan)

	return errorsChan
}

func (bot *Bot) startProcessing(errorsChan chan error) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.BotAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			command := update.Message.Text
			command = strings.Trim(command, " ")
			if command[:1] == "/" {

				user := storage.User{
					Name: update.Message.From.UserName,
					TGId: update.Message.From.ID,
				}
				bot.addUserIfNotExists(&user)

				msgText := bot.commandProcessing(command, &user)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ReplyToMessageID = update.Message.MessageID

				_, err := bot.BotAPI.Send(msg)
				if err != nil {
					log.Println("Error in Dial", err)
					errorsChan <- err
				}
			}
		}
	}
}

//commandProcessing - Processing known commands
func (bot *Bot) commandProcessing(command string, user *storage.User) string {
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
		return "/help - print help message\n" +
			"/check www.checkURL1.com www.checkURL2.com ... - check certificate on URL. Use spaces to check few domains"
	case "/check":
		if attr == "" {
			return "You must specify the URL. Format: \n\t /check www.checkURL1.com www.checkURL2.com ... Use space to check few URLs."
		}
		return certinfo.GetCertsInfo(attr, false)
	default:
		return "Use /help command"
	}
}

//addUserIfNotExists - add new user to storage
func (bot *Bot) addUserIfNotExists(user *storage.User) {
	savedUser, err := bot.db.GetUserByTGId(user.TGId)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotFound) {
			_, err := bot.db.AddUser(user)
			if err != nil {
				log.Println(err)
				return
			}
			return
		}
		log.Println(err)
		return
	}
	if user.Name != savedUser.Name {
		savedUser.Name = user.Name
		_, err2 := bot.db.UpdateUserInfo(savedUser)
		if err2 != nil {
			log.Println(err2)
			return
		}
	}
	user = savedUser
	return
}
