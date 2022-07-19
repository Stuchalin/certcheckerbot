package sqlite3

import (
	storage "certcheckerbot/storage"
	"fmt"
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
			defer got.Dispose()
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
			defer controller.Dispose()

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

func TestSqlite3Controller_GetUserByTGId(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	type args struct {
		tgId string
	}
	tests := []struct {
		name    string
		args    args
		want    *storage.User
		wantErr bool
	}{
		{
			name: "test GetUserByTGId",
			args: args{
				tgId: "test",
			},
			want: &storage.User{
				Name:             "test",
				TGId:             "test",
				NotificationHour: 0,
				UTC:              0,
			},
			wantErr: false,
		},
		{
			name: "test GetUserByTGId - fail to find",
			args: args{
				tgId: "test2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewController(dbName)
			defer db.Dispose()
			if tt.want != nil {
				tt.want.Id, _ = db.AddUser(tt.want)
			}

			got, err := db.GetUserByTGId(tt.args.tgId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByTGId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByTGId() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqlite3Controller_GetUserName(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *storage.User
		wantErr bool
	}{
		{
			name: "test GetUserByName",
			args: args{
				name: "test",
			},
			want: &storage.User{
				Name:             "test",
				TGId:             "test",
				NotificationHour: 0,
				UTC:              0,
			},
			wantErr: false,
		},
		{
			name: "test GetUserByName - fail to find",
			args: args{
				name: "test2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewController(dbName)
			defer db.Dispose()
			if tt.want != nil {
				tt.want.Id, _ = db.AddUser(tt.want)
			}

			got, err := db.GetUserByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqlite3Controller_AddUserDomain_FK_error(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	type args struct {
		domain *storage.UserDomain
	}
	tests := []struct {
		name              string
		args              args
		want              bool
		wantErr           bool
		expectedErrorText string
	}{
		{
			name: "test AddUserDomain FK Error",
			args: args{&storage.UserDomain{
				UserId: 1,
				Domain: "test.com",
			}},
			want:              false,
			wantErr:           true,
			expectedErrorText: "FOREIGN KEY constraint failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewController(dbName)
			defer db.Dispose()

			got, err := db.AddUserDomain(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddUserDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddUserDomain() got = %v, want %v", got, tt.want)
				return
			}
			if err.Error() != tt.expectedErrorText {
				t.Errorf("AddUserDomain() got = %v, want %v", err, tt.expectedErrorText)
				return
			}
		})
	}
}

func TestSqlite3Controller_AddUserDomain(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	type args struct {
		domain *storage.UserDomain
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test AddUserDomain correct",
			args: args{&storage.UserDomain{
				UserId: 1,
				Domain: "test.com",
			}},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewController(dbName)
			defer db.Dispose()

			user := storage.User{
				Name: "test",
				TGId: "test",
			}

			_, _ = db.AddUser(&user)
			tt.args.domain.UserId = user.Id

			got, err := db.AddUserDomain(tt.args.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddUserDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddUserDomain() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestSqlite3Controller_GetUserDomains(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		{
			name:    "test GetUserDomains correct",
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewController(dbName)
			defer db.Dispose()

			user := storage.User{
				Name: "test",
				TGId: "test",
			}

			_, _ = db.AddUser(&user)

			domain := storage.UserDomain{
				UserId: user.Id,
				Domain: "test.com",
			}
			domain2 := storage.UserDomain{
				UserId: user.Id,
				Domain: "test2.ru",
			}

			var domains []storage.UserDomain

			domains = append(domains, domain)
			domains = append(domains, domain2)

			for _, dom := range domains {
				_, _ = db.AddUserDomain(&dom)
			}

			result, err := db.GetUserDomains(&user)
			if err != nil {
				t.Errorf("GetUserDomains() return error %v", err)
				return
			}
			if !reflect.DeepEqual(result, &domains) {
				t.Errorf("GetUserDomains() got %v, want %v", result, domains)
				return
			}
		})
	}
}

func TestSqlite3Controller_GetUserDomains_no_domains(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	tests := []struct {
		name string
		want string
	}{
		{
			name: "test GetUserDomains no domains",
			want: "cannot find user domains for user test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewController(dbName)
			defer db.Dispose()

			user := storage.User{
				Name: "test",
				TGId: "test",
			}

			_, err := db.GetUserDomains(&user)
			if err == nil {
				t.Errorf("GetUserDomains() expected error  %s", tt.want)
				return
			}
			if err.Error() != tt.want {
				t.Errorf("GetUserDomains()  got %v, want %v", err, tt.want)
				return
			}
		})
	}
}

func TestSqlite3Controller_RemoveUserDomain(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		{
			name:    "test RemoveUserDomain success",
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewController(dbName)
			defer db.Dispose()

			user := storage.User{
				Name: "test",
				TGId: "test",
			}

			_, _ = db.AddUser(&user)

			domain := storage.UserDomain{
				UserId: user.Id,
				Domain: "test.com",
			}
			domain2 := storage.UserDomain{
				UserId: user.Id,
				Domain: "test2.ru",
			}

			var domains []storage.UserDomain

			domains = append(domains, domain)
			domains = append(domains, domain2)

			for _, dom := range domains {
				_, _ = db.AddUserDomain(&dom)
			}

			domains = domains[1:]

			got, err := db.RemoveUserDomain(&domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveUserDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RemoveUserDomain() got = %v, want %v", got, tt.want)
				return
			}
			result, err := db.GetUserDomains(&user)
			if err != nil {
				t.Errorf("RemoveUserDomain() return error %v", err)
				return
			}
			if !reflect.DeepEqual(result, &domains) {
				t.Errorf("RemoveUserDomain() got %v, want %v", result, domains)
				return
			}
		})
	}
}

func Test_removeAllUserDomains(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	tests := []struct {
		name       string
		want       bool
		wantErr    bool
		errMessage string
	}{
		{
			name:       "test removeAllUserDomains success",
			want:       true,
			wantErr:    false,
			errMessage: "cannot find user domains for user test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewController(dbName)
			defer db.Dispose()

			user := storage.User{
				Name: "test",
				TGId: "test",
			}

			_, _ = db.AddUser(&user)

			domain := storage.UserDomain{
				UserId: user.Id,
				Domain: "test.com",
			}
			domain2 := storage.UserDomain{
				UserId: user.Id,
				Domain: "test2.ru",
			}

			var domains []storage.UserDomain

			domains = append(domains, domain)
			domains = append(domains, domain2)

			for _, dom := range domains {
				_, _ = db.AddUserDomain(&dom)
			}

			tx, _ := db.Connection.Begin()
			got, err := removeAllUserDomains(&user, tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("removeAllUserDomains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = tx.Commit()
			if got != tt.want {
				t.Errorf("removeAllUserDomains() got = %v, want %v", got, tt.want)
				return
			}
			_, err = db.GetUserDomains(&user)
			if err == nil {
				t.Errorf("removeAllUserDomains() expected error %s", tt.errMessage)
				return
			}
			if err.Error() != tt.errMessage {
				t.Errorf("removeAllUserDomains() got = %v, want %s", err, tt.errMessage)
				return
			}
		})
	}
}

func TestSqlite3Controller_RemoveUser(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	tests := []struct {
		name              string
		want              bool
		wantErr           bool
		errMessageDomains string
		errMessageUser    string
	}{
		{
			name:              "test RemoveUser success",
			want:              true,
			wantErr:           false,
			errMessageDomains: "cannot find user domains for user test",
			errMessageUser:    "cannot find user by id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewController(dbName)
			defer db.Dispose()

			user := storage.User{
				Name: "test",
				TGId: "test",
			}

			_, _ = db.AddUser(&user)

			domain := storage.UserDomain{
				UserId: user.Id,
				Domain: "test.com",
			}
			domain2 := storage.UserDomain{
				UserId: user.Id,
				Domain: "test2.ru",
			}

			var domains []storage.UserDomain

			domains = append(domains, domain)
			domains = append(domains, domain2)

			for _, dom := range domains {
				_, _ = db.AddUserDomain(&dom)
			}

			got, err := db.RemoveUser(&user)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RemoveUser() got = %v, want %v", got, tt.want)
				return
			}
			_, err = db.GetUserDomains(&user)
			if err == nil {
				t.Errorf("RemoveUser() expected error %s", tt.errMessageDomains)
				return
			}
			if err.Error() != tt.errMessageDomains {
				t.Errorf("RemoveUser() got = %v, want %s", err, tt.errMessageDomains)
				return
			}
			_, err = db.GetUserById(user.Id)
			tt.errMessageUser = tt.errMessageUser + " " + fmt.Sprint(user.Id)
			if err == nil {
				t.Errorf("RemoveUser() expected error %s", tt.errMessageUser)
				return
			}
			if err.Error() != tt.errMessageUser {
				t.Errorf("RemoveUser() got = %v, want %s", err, tt.errMessageUser)
				return
			}
		})
	}
}

func TestSqlite3Controller_UpdateUserInfo(t *testing.T) {
	dbName := getTempDBName()
	defer removeDbFile(dbName)

	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		{
			name:    "test UpdateUserInfo success",
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _ := NewController(dbName)
			defer db.Dispose()

			userStart := storage.User{
				Name: "test",
				TGId: "test",
			}

			userChanged := storage.User{
				Name:             "test new",
				TGId:             "test new",
				NotificationHour: 10,
				UTC:              10,
			}

			_, _ = db.AddUser(&userStart)
			userChanged.Id = userStart.Id

			got, err := db.UpdateUserInfo(&userChanged)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateUserInfo() got = %v, want %v", got, tt.want)
				return
			}
			dbUser, err := db.GetUserById(userStart.Id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(dbUser, &userChanged) {
				t.Errorf("UpdateUserInfo() got %v, want %v", dbUser, userChanged)
				return
			}
		})
	}
}
