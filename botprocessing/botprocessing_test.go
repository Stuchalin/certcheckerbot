package botprocessing

import (
	"certcheckerbot/storage"
	"certcheckerbot/storage/sqlite3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"os"
	"regexp"
	"testing"
)

func getTempDBName() string {
	return os.TempDir() + uuid.New().String() + ".db"
}

func removeDbFile(dbName string) {
	_ = os.Remove(dbName)
}

func TestBot_commandProcessing(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	db, _ := sqlite3.NewController(dbName)
	defer db.Dispose()

	user := storage.User{
		Id:               1,
		Name:             "test user",
		TGId:             123,
		NotificationHour: 0,
		UTC:              0,
	}
	domain := storage.UserDomain{
		UserId: 1,
		Domain: "google.com",
	}
	_, _ = db.AddUser(&user)
	_, _ = db.AddUserDomain(&domain)

	userForSelectNoDomains := storage.User{
		Id:               2,
		Name:             "test user 2",
		TGId:             12345,
		NotificationHour: 0,
		UTC:              0,
	}
	_, _ = db.AddUser(&userForSelectNoDomains)

	userForSelectExistDomain := storage.User{
		Id:               3,
		Name:             "test user 3",
		TGId:             1234567,
		NotificationHour: 0,
		UTC:              0,
	}
	domainForSelectExistDomain := storage.UserDomain{
		UserId: 3,
		Domain: "google.com",
	}
	_, _ = db.AddUser(&userForSelectExistDomain)
	_, _ = db.AddUserDomain(&domainForSelectExistDomain)

	userForRemoveDomain := storage.User{
		Id:               4,
		Name:             "test user 4",
		TGId:             12345678,
		NotificationHour: 0,
		UTC:              0,
	}
	domainForRemoveDomain := storage.UserDomain{
		UserId: 4,
		Domain: "google.com",
	}
	_, _ = db.AddUser(&userForRemoveDomain)
	_, _ = db.AddUserDomain(&domainForRemoveDomain)

	type fields struct {
		BotAPI *tgbotapi.BotAPI
		db     storage.UsersConfig
	}
	type args struct {
		command string
		user    *storage.User
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      string
		wantRegex string
	}{
		//help command
		{
			name:   "test /help",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/help",
			},
			want: "Simple bot for check certificates expire dates\n" +
				"version 0.1\n" +
				"\t/help - print help message\n" +
				"\t/check www.checkURL1.com www.checkURL2.com ... - check certificate on URL. Use spaces to check few domains\n" +
				"\t/set_hour [hour in 24 format 0..23] - set a notification hour for messages about expired domains. For example: \"/set_hour 9\". Notification hour for default - 0.\n" +
				"\t/set_tz [-11..14] - set a timezone for messages about expired domains. For example: \\\"/set_tz 3\\\". Timezone for default - 0.\n" +
				"\t/domains - get added domains\n" +
				"\t/add_domain [domain_name] - add domain for schedule checks. For example: \"/add_domain google.com\"\n" +
				"\t/remove_domain [domain_name] - removes domain for schedule checks. For example: \"/remove_domain google.com\"\n",
		},
		//empty command
		{
			name:   "test empty command",
			fields: fields{},
			args: args{
				user:    nil,
				command: "",
			},
			want: "Use /help command",
		},
		//unknown command
		{
			name:   "test unknown command",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/12315asg",
			},
			want: "Use /help command",
		},
		//set_hour
		{
			name:   "test /set_hour with no attrs",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour",
			},
			want: "You must specify the notification hour. Format: \n\t /set_hour [hour in 24 format 0..23]. For example: \"/set_hour 9\"",
		},
		{
			name:   "test /set_hour not int hour",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour qwe",
			},
			want: "Notification hour must be integer number in 0..23 range.",
		},
		{
			name:   "test /set_hour decimal hour",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour 1.1",
			},
			want: "Notification hour must be integer number in 0..23 range.",
		},
		{
			name:   "test /set_hour not correct hour -2",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour -2",
			},
			want: "Notification hour must be integer number in 0..23 range.",
		},
		{
			name:   "test /set_hour not correct hour 24",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour 24",
			},
			want: "Notification hour must be integer number in 0..23 range.",
		},
		{
			name:   "test /set_hour user not identified",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour 23",
			},
			want: "Internal error: user not identified",
		},
		{
			name:   "test /set_hour no user for update",
			fields: fields{db: db},
			args: args{
				user: &storage.User{
					Id:               -100,
					Name:             "test user",
					TGId:             123,
					NotificationHour: 0,
					UTC:              0,
				},
				command: "/set_hour 23",
			},
			want: "Internal error: cannot update user notification hour",
		},
		{
			name:   "test /set_hour success test",
			fields: fields{db: db},
			args: args{
				user:    &user,
				command: "/set_hour 23",
			},
			want: "Notification hour is successful set on 23",
		},
		//set_tz
		{
			name:   "test /set_tz with no attrs",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_tz",
			},
			want: "You must specify the timezone. Format: \n\t /set_tz [-11..14]. For example: \"/set_tz 3\"",
		},
		{
			name:   "test /set_tz not int tz",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_tz qwe",
			},
			want: "Timezone must be integer number in -11..14 range.",
		},
		{
			name:   "test /set_tz decimal tz",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_tz 1.1",
			},
			want: "Timezone must be integer number in -11..14 range.",
		},
		{
			name:   "test /set_tz not correct tz -12",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_tz -12",
			},
			want: "Timezone must be integer number in -11..14 range.",
		},
		{
			name:   "test /set_tz not correct tz 15",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_tz 15",
			},
			want: "Timezone must be integer number in -11..14 range.",
		},
		{
			name:   "test /set_tz user not identified",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_tz 3",
			},
			want: "Internal error: user not identified",
		},
		{
			name:   "test /set_tz no user for update",
			fields: fields{db: db},
			args: args{
				user: &storage.User{
					Id:               -100,
					Name:             "test user",
					TGId:             123,
					NotificationHour: 0,
					UTC:              0,
				},
				command: "/set_tz 3",
			},
			want: "Internal error: cannot update timezone",
		},
		{
			name:   "test /set_tz success test",
			fields: fields{db: db},
			args: args{
				user:    &user,
				command: "/set_tz 3",
			},
			want: "Timezone is successful set on 3",
		},
		{
			name:   "test /set_tz success test with plus",
			fields: fields{db: db},
			args: args{
				user:    &user,
				command: "/set_tz +3",
			},
			want: "Timezone is successful set on 3",
		},
		{
			name:   "test /set_tz success test with minus",
			fields: fields{db: db},
			args: args{
				user:    &user,
				command: "/set_tz -11",
			},
			want: "Timezone is successful set on -11",
		},
		{
			name:   "test /add_domain with no attrs",
			fields: fields{db: db},
			args: args{
				user:    &user,
				command: "/add_domain ",
			},
			want: "You must specify domain name. Format: \n\t /add_domain [domain_name]. For example: \"/add_domain google.com\"",
		},
		{
			name:   "test /add_domain cannot add multiply domains error",
			fields: fields{db: db},
			args: args{
				user:    &user,
				command: "/add_domain google.com twitch.com",
			},
			want: "You cannot add multiple domains at once. Please specify only one domain.",
		},
		{
			name:   "test /add_domain no such host error",
			fields: fields{db: db},
			args: args{
				user:    &user,
				command: "/add_domain www",
			},
			wantRegex: "Fail add domain for schedule checks. \nCannot check certificate for this domain. Error: check certificate error - cannot check cert from URL www. Error: .*no such host www.*",
		},
		{
			name:   "test /add_domain domain already added",
			fields: fields{db: db},
			args: args{
				user:    &user,
				command: "/add_domain google.com",
			},
			want: "Fail add domain - google.com. This domain already added to account. Check added domains with command /domains",
		},
		{
			name:   "test /add_domain success add domain",
			fields: fields{db: db},
			args: args{
				user:    &user,
				command: "/add_domain ya.ru",
			},
			want: "Domain successfully added.",
		},
		{
			name:   "test /domains no domains",
			fields: fields{db: db},
			args: args{
				user:    &userForSelectNoDomains,
				command: "/domains",
			},
			want: "You have no added domains.",
		},
		{
			name:   "test /domains success get domains",
			fields: fields{db: db},
			args: args{
				user:    &userForSelectExistDomain,
				command: "/domains",
			},
			want: "Added domains:\n\tgoogle.com\n",
		},
		{
			name:   "test /remove_domain with no attrs",
			fields: fields{db: db},
			args: args{
				user:    &userForRemoveDomain,
				command: "/remove_domain ",
			},
			want: "You must specify domain name. Format: \n\t /remove_domain [domain_name]. For example: \"/remove_domain google.com\"",
		},
		{
			name:   "test /remove_domain cannot add multiply domains error",
			fields: fields{db: db},
			args: args{
				user:    &userForRemoveDomain,
				command: "/remove_domain google.com twitch.com",
			},
			want: "You cannot remove multiple domains at once. Please specify only one domain.",
		},
		{
			name:   "test /remove_domain this domain does not added for you",
			fields: fields{db: db},
			args: args{
				user:    &userForRemoveDomain,
				command: "/remove_domain twitch.com",
			},
			want: "Fail to remove domain, this domain does not added for you. To check added domains use /domains command.",
		},

		{
			name:   "test /remove_domain success remove domain",
			fields: fields{db: db},
			args: args{
				user:    &userForRemoveDomain,
				command: "/remove_domain google.com",
			},
			want: "Domain successfully removed.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot{
				BotAPI: tt.fields.BotAPI,
				db:     tt.fields.db,
			}
			if tt.wantRegex != "" {
				got := bot.commandProcessing(tt.args.command, tt.args.user)
				res, err := regexp.MatchString(tt.wantRegex, got)
				if err != nil {
					t.Errorf("commandProcessing() - regex error: %s", err)
				}
				if !res {
					t.Errorf("commandProcessing() = %v, regex pattern = %v", got, tt.wantRegex)
				}
			} else if got := bot.commandProcessing(tt.args.command, tt.args.user); got != tt.want {
				t.Errorf("commandProcessing() = %v, want %v", got, tt.want)
			}
		})
	}
}
