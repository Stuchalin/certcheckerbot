package sqlite3

import (
	storage "certcheckerbot/storage"
	"github.com/google/uuid"
	"os"
	"reflect"
	"testing"
)

func getTempDBName() string {
	return os.TempDir() + uuid.New().String() + ".db"
}

func removeDbFile(dbName string) {
	_ = os.Remove(dbName)
}

func TestNewController(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	type args struct {
		databaseName string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test TestNewController created",
			args: args{
				databaseName: dbName,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewController(tt.args.databaseName)
			defer DisposeConnection(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewController() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				version, err := GetCurrentDBVersion(got.Connection)
				if err != nil {
					t.Errorf("NewController() cannot get version from database (%v)", err)
					return
				}
				if version >= 1 {
					return
				} else {
					t.Errorf("Incorrect version of database %d", version)
				}
			}
		})
	}
}

func TestSqlite3Controller_AddUser(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	type args struct {
		user *storage.User
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test Sqlite3Controller_AddUser correct",
			args: args{user: &storage.User{
				Name:             "test",
				TGId:             "test",
				NotificationHour: 1,
				UTC:              1,
			}},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, _ := NewController(dbName)
			defer DisposeConnection(controller)

			got, err := controller.AddUser(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddUser() got = %v, want %v", got, tt.want)
				return
			}
			if tt.args.user.Id != tt.want {
				t.Errorf("AddUser() tt.args.user.Id = %v, want %v", got, tt.want)
				return
			}
			user, err := controller.GetUserById(tt.args.user.Id)
			if err != nil {
				t.Errorf("AddUser() cannot find user after add by id %d", tt.args.user.Id)
				return
			}
			if !reflect.DeepEqual(tt.args.user, user) {
				t.Errorf("AddUser() selectd user afrer add in not equals with adding. Added = %v, selected = %v", tt.args.user, user)
			}
		})
	}
}
