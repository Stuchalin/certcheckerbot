package sqlite3

import (
	"database/sql"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init_db() (string, *sql.DB) {
	dbName := os.TempDir() + uuid.New().String() + ".db"
	file, _ := os.Create(dbName)
	file.Close()

	db, _ := sql.Open("sqlite3", dbName)

	return dbName, db
}

func remove_db(dbName string, db *sql.DB) {
	db.Close()
	os.Remove(dbName)
}

func Test_isBaseStructExists_NoBaseStruct(t *testing.T) {
	dbName, db := init_db()
	defer remove_db(dbName, db)

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "test no base struct",
			args:    args{db: db},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isBaseStructExists(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("isBaseStructExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isBaseStructExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isBaseStructExists_NoDatabase(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name:    "test base struct nil database",
			args:    args{db: nil},
			wantErr: "Nil database connection",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := isBaseStructExists(tt.args.db)
			if err == nil {
				t.Errorf("isBaseStructExists() expected error %v", tt.wantErr)
				return
			} else if err.Error() != tt.wantErr {
				t.Errorf("isBaseStructExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_isBaseStructExists_ClosedDatabase(t *testing.T) {
	dbName, db := init_db()
	db.Close()
	defer remove_db(dbName, db)

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name:    "test base struct closed database",
			args:    args{db: db},
			wantErr: "sql: database is closed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := isBaseStructExists(tt.args.db)
			if err == nil {
				t.Errorf("isBaseStructExists() expected error %v", tt.wantErr)
				return
			} else if err.Error() != tt.wantErr {
				t.Errorf("isBaseStructExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_isBaseStructExists_ExistBaseStruct(t *testing.T) {
	dbName, db := init_db()
	tx, _ := db.Begin()
	tx.Exec("create table db_versions (version integer)")
	tx.Commit()
	defer remove_db(dbName, db)

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "test exist base struct",
			args:    args{db: db},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isBaseStructExists(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("isBaseStructExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isBaseStructExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentDBVersion_NoStruct(t *testing.T) {
	dbName, db := init_db()
	defer remove_db(dbName, db)

	type args struct {
		db *sql.DB
	}
	var tests = []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "test no database struct",
			args:    args{db},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCurrentDBVersion(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentDBVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCurrentDBVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentDBVersion_NoVersions(t *testing.T) {
	dbName, db := init_db()
	tx, _ := db.Begin()
	tx.Exec("create table db_versions (version integer)")
	tx.Commit()
	defer remove_db(dbName, db)

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "test no versions in db",
			args:    args{db},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCurrentDBVersion(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentDBVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCurrentDBVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentDBVersion_GetMaxVersion(t *testing.T) {
	dbName, db := init_db()
	tx, _ := db.Begin()
	tx.Exec("create table db_versions (version integer) ")
	insertStatement, _ := tx.Prepare("insert into db_versions(version) values (?)")
	maxVersion := 0
	for i := 1; i <= 15; i++ {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		version := r1.Intn(100)
		insertStatement.Exec(version)
		if maxVersion < version {
			maxVersion = version
		}
	}
	tx.Commit()
	defer remove_db(dbName, db)

	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "test get max version",
			args:    args{db},
			want:    maxVersion,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCurrentDBVersion(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentDBVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetCurrentDBVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
