package botprocessing

import (
	"certcheckerbot/storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"testing"
)

func TestBot_commandProcessing(t *testing.T) {
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
		{
			name:   "test /help",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/help",
			},
			want: "/help - print help message\n" +
				"/check www.checkURL1.com www.checkURL2.com ... - check certificate on URL. Use spaces to check few domains",
		},
		{
			name:   "test empty command",
			fields: fields{},
			args: args{
				user:    nil,
				command: "",
			},
			want: "Use /help command",
		},
		{
			name:   "test unknown command",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/12315asg",
			},
			want: "Use /help command",
		},
		{
			name:   "test /set_hour with no attrs",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour",
			},
			want: "You must specify the reminder hour. Format: \n\t /set_hour [hour in 24 format 0..23]. For example: \"/set_hour 9\"",
		},
		{
			name:   "test /set_hour not int hour",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour qwe",
			},
			want: "Reminder hour must be integer number in 0..23 range.",
		},
		{
			name:   "test /set_hour decimal hour",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour 1.1",
			},
			want: "Reminder hour must be integer number in 0..23 range.",
		},
		{
			name:   "test /set_hour not correct hour -2",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour -2",
			},
			want: "Reminder hour must be integer number in 0..23 range.",
		},
		{
			name:   "test /set_hour not correct hour 24",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour 24",
			},
			want: "Reminder hour must be integer number in 0..23 range.",
		},
		{
			name:   "test /set_hour success test",
			fields: fields{},
			args: args{
				user:    nil,
				command: "/set_hour 23",
			},
			want: "Reminder hour is successful set on 23",
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
