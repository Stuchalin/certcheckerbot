package botprocessing

import (
	"certcheckerbot/storage"
	"certcheckerbot/storage/sqlite3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"os"
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
	_, _ = db.AddUser(&user)

	type fields struct {
		BotAPI *tgbotapi.BotAPI
		db     storage.UsersConfig
	}
	type args struct {
		command string
		user    *storage.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
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
				"\t/set_tz [-11..14] - set a timezone for messages about expired domains. For example: \\\"/set_tz 3\\\". Timezone for default - 0.",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot := &Bot{
				BotAPI: tt.fields.BotAPI,
				db:     tt.fields.db,
			}
			if got := bot.commandProcessing(tt.args.command, tt.args.user); got != tt.want {
				t.Errorf("commandProcessing() = %v, want %v", got, tt.want)
			}
		})
	}
}
