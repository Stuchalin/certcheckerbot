package botprocessing

import (
	"certcheckerbot/certinfo"
	"certcheckerbot/storage"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
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

				user := &storage.User{
					Name: update.Message.From.UserName,
					TGId: update.Message.From.ID,
				}
				user = bot.addUserIfNotExists(user)

				msgText := bot.commandProcessing(command, user)

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
		return "Simple bot for check certificates expire dates\n" +
			"version 0.2\n" +
			"\t/help - print help message\n" +
			"\t/check www.checkURL1.com www.checkURL2.com ... - check certificate on URL. Use spaces to check few domains\n" +
			"\t/set_hour [hour in 24 format 0..23] - set a notification hour for messages about expired domains. For example: \"/set_hour 9\". Notification hour for default - 0.\n" +
			"\t/set_tz [-11..14] - set a timezone for messages about expired domains. For example: \\\"/set_tz 3\\\". Timezone for default - 0.\n" +
			"\t/domains - get added domains\n" +
			"\t/add_domain [domain_name] - add domain for schedule checks. For example: \"/add_domain google.com\"\n" +
			"\t/remove_domain [domain_name] - removes domain for schedule checks. For example: \"/remove_domain google.com\"\n"
	case "/check":
		if attr == "" {
			return "You must specify the URL. Format: \n\t /check www.checkURL1.com www.checkURL2.com ... Use space to check few URLs."
		}
		return certinfo.GetCertsInfo(attr, false)
	case "/set_hour":
		if attr == "" {
			return "You must specify the notification hour. Format: \n\t /set_hour [hour in 24 format 0..23]. For example: \"/set_hour 9\""
		}
		hour, err := strconv.Atoi(attr)
		if err != nil || hour < 0 || hour > 23 {
			return "Notification hour must be integer number in 0..23 range."
		}
		if user == nil {
			log.Println("Internal error: user not identified")
			return "Internal error: user not identified"
		}
		user.NotificationHour = hour
		result, err := bot.db.UpdateUserInfo(user)
		if err != nil {
			log.Println(err)
			return fmt.Sprintf("Internal error: cannot set notification hour to %s", attr)
		}
		if result {
			return fmt.Sprintf("Notification hour is successful set on %s", attr)
		} else {
			log.Println("Internal error: cannot update user notification hour")
			return "Internal error: cannot update user notification hour"
		}
	case "/set_tz":
		if attr == "" {
			return "You must specify the timezone. Format: \n\t /set_tz [-11..14]. For example: \"/set_tz 3\""
		}
		attr = strings.Replace(attr, "+", "", -1)
		tz, err := strconv.Atoi(attr)
		if err != nil || tz < -11 || tz > 14 {
			return "Timezone must be integer number in -11..14 range."
		}
		if user == nil {
			log.Println("Internal error: user not identified")
			return "Internal error: user not identified"
		}
		user.UTC = tz
		result, err := bot.db.UpdateUserInfo(user)
		if err != nil {
			log.Println(err)
			return fmt.Sprintf("Internal error: cannot set timezone to %s", attr)
		}
		if result {
			return fmt.Sprintf("Timezone is successful set on %s", attr)
		} else {
			log.Println("Internal error: cannot update timezone")
			return "Internal error: cannot update timezone"
		}

	case "/add_domain":
		if attr == "" {
			return "You must specify domain name. Format: \n\t /add_domain [domain_name]. For example: \"/add_domain google.com\""
		}

		if strings.Contains(attr, " ") {
			return "You cannot add multiple domains at once. Please specify only one domain."
		}

		_, err := certinfo.GetCertInfo(attr, false)
		if err != nil {
			return fmt.Sprintf("Fail add domain for schedule checks. \nCannot check certificate for this domain. Error: %v", err)
		}

		newUserDomain := storage.UserDomain{UserId: user.Id, Domain: attr}

		result, err := bot.db.AddUserDomain(&newUserDomain)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				return fmt.Sprintf("Fail add domain - %s. This domain already added to account. Check added domains with command /domains", attr)
			}
			log.Println(fmt.Sprintf("Internal error: Fail to add domain. Error: %v.", err))
			return fmt.Sprintf("Internal error: Fail to add domain. Error: %v.", err)
		}
		if !result {
			log.Println("Internal error: Fail to add domain.")
			return fmt.Sprintf("Internal error: Fail to add domain.")
		}

		return "Domain successfully added."

	case "/domains":
		domains, err := bot.db.GetUserDomains(user)
		if err != nil {
			if err == storage.ErrorUserDomainNotFound {
				return fmt.Sprintf("You have no added domains.")
			}
			log.Println(fmt.Sprintf("Internal error: cannot get user domains. Error: %v", err))
			return fmt.Sprintf("Internal error: cannot get user domains. Error: %v", err)
		}

		domainsResult := "Added domains:\n"

		for _, domain := range *domains {
			domainsResult += "\t" + domain.Domain + "\n"
		}

		return domainsResult

	case "/remove_domain":
		if attr == "" {
			return "You must specify domain name. Format: \n\t /remove_domain [domain_name]. For example: \"/remove_domain google.com\""
		}

		if strings.Contains(attr, " ") {
			return "You cannot remove multiple domains at once. Please specify only one domain."
		}

		newUserDomain := storage.UserDomain{UserId: user.Id, Domain: attr}

		result, err := bot.db.RemoveUserDomain(&newUserDomain)
		if err != nil {
			log.Println(fmt.Sprintf("Internal error: Fail to remove domain. Error: %v.", err))
			return fmt.Sprintf("Internal error: Fail to remove domain. Error: %v.", err)
		}
		if !result {
			return fmt.Sprintf("Fail to remove domain, this domain does not added for you. To check added domains use /domains command.")
		}

		return "Domain successfully removed."

	default:
		return "Use /help command"
	}
}

//addUserIfNotExists - add new user to storage
func (bot *Bot) addUserIfNotExists(user *storage.User) *storage.User {
	savedUser, err := bot.db.GetUserByTGId(user.TGId)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotFound) {
			user.Id, err = bot.db.AddUser(user)
			if err != nil {
				log.Println(err)
				return nil
			}
			return user
		}
		log.Println(err)
		return nil
	}
	if user.Name != savedUser.Name {
		savedUser.Name = user.Name
		_, err2 := bot.db.UpdateUserInfo(savedUser)
		if err2 != nil {
			log.Println(err2)
			return nil
		}
	}
	user = savedUser
	return user
}
