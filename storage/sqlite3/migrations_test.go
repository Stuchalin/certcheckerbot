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
	_ = file.Close()

	db, _ := sql.Open("sqlite3", dbName)

	return dbName, db
}

func remove_db(dbName string, db *sql.DB) {
	_ = db.Close()
	_ = os.Remove(dbName)
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
	tx.Exec("create table db_versions (version integer, PRIMARY KEY (version));")
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
	tx.Exec("create table db_versions (version integer, PRIMARY KEY (version));")
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
	tx.Exec("create table db_versions (version integer, PRIMARY KEY (version));")
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

func Test_setDBVersion_correct_nilTran(t *testing.T) {
	dbName, db := init_db()
	tx, _ := db.Begin()
	tx.Exec("create table db_versions (version integer, PRIMARY KEY (version));")
	tx.Commit()
	defer remove_db(dbName, db)

	type args struct {
		version int
		db      *sql.DB
		tx      *sql.Tx
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test correct set db version - nil tran",
			args: args{
				version: 1,
				db:      db,
				tx:      nil,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setDBVersion(tt.args.version, tt.args.db, tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("setDBVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("setDBVersion() got = %v, want %v", got, tt.want)
				return
			}
			rows, _ := db.Query("select version from db_versions where version = ?", tt.args.version)
			defer rows.Close()

			version := -1
			for rows.Next() {
				rows.Scan(&version)
			}
			if version != tt.args.version {
				t.Errorf("setDBVersion() getVersion = %v, setVersion %v", version, tt.args.version)
			}
		})
	}
}

func Test_setDBVersion_correct_sendTran_commit(t *testing.T) {
	dbName, db := init_db()
	tx, _ := db.Begin()
	tx.Exec("create table db_versions (version integer, PRIMARY KEY (version));")
	tx.Commit()
	defer remove_db(dbName, db)

	type args struct {
		version int
		db      *sql.DB
		tx      *sql.Tx
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test correct set db version - send tran commit",
			args: args{
				version: 1,
				db:      db,
				tx:      nil,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, _ := tt.args.db.Begin()
			got, err := setDBVersion(tt.args.version, tt.args.db, tx)
			tx.Commit()
			if (err != nil) != tt.wantErr {
				t.Errorf("setDBVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("setDBVersion() got = %v, want %v", got, tt.want)
				return
			}
			rows, _ := db.Query("select version from db_versions where version = ?", tt.args.version)
			defer rows.Close()

			version := -1
			for rows.Next() {
				rows.Scan(&version)
			}
			if version != tt.args.version {
				t.Errorf("setDBVersion() getVersion = %v, setVersion %v", version, tt.args.version)
			}
		})
	}
}

func Test_setDBVersion_correct_sendTran_rollback(t *testing.T) {
	dbName, db := init_db()
	tx, _ := db.Begin()
	tx.Exec("create table db_versions (version integer, PRIMARY KEY (version));")
	tx.Commit()
	defer remove_db(dbName, db)

	type args struct {
		version int
		db      *sql.DB
		tx      *sql.Tx
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test correct set db version - send tran rollback",
			args: args{
				version: 1,
				db:      db,
				tx:      nil,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, _ := tt.args.db.Begin()
			got, err := setDBVersion(tt.args.version, tt.args.db, tx)
			tx.Rollback()
			if (err != nil) != tt.wantErr {
				t.Errorf("setDBVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("setDBVersion() got = %v, want %v", got, tt.want)
				return
			}
			rows, _ := db.Query("select version from db_versions where version = ?", tt.args.version)
			defer rows.Close()

			version := -1
			for rows.Next() {
				rows.Scan(&version)
			}
			if version == tt.args.version {
				t.Errorf("setDBVersion() getVersion = %v, setVersion %v", version, tt.args.version)
			}
		})
	}
}

func Test_setDBVersion_add_existing_version_nilTran(t *testing.T) {
	dbName, db := init_db()
	tx, _ := db.Begin()
	tx.Exec("create table db_versions (version integer, PRIMARY KEY (version));")
	tx.Commit()
	defer remove_db(dbName, db)

	type args struct {
		version int
		db      *sql.DB
		tx      *sql.Tx
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test add existing version - nil tran",
			args: args{
				version: 1,
				db:      db,
				tx:      nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := setDBVersion(tt.args.version, tt.args.db, tt.args.tx)
			got, err = setDBVersion(tt.args.version, tt.args.db, tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("setDBVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("setDBVersion() got = %v, want %v", got, tt.want)
				return
			}
			rows, _ := db.Query("select version from db_versions where version = ?", tt.args.version)
			defer rows.Close()

			version := -1
			for rows.Next() {
				rows.Scan(&version)
			}
			if version != tt.args.version {
				t.Errorf("setDBVersion() getVersion = %v, setVersion %v", version, tt.args.version)
			}
		})
	}
}

func Test_setDBVersion_add_existing_version_sendTran_rollback(t *testing.T) {
	dbName, db := init_db()
	tx, _ := db.Begin()
	tx.Exec("create table db_versions (version integer, PRIMARY KEY (version));")
	tx.Commit()
	defer remove_db(dbName, db)

	type args struct {
		version int
		db      *sql.DB
		tx      *sql.Tx
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test add existing version - send tran rollback",
			args: args{
				version: 1,
				db:      db,
				tx:      nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, _ := tt.args.db.Begin()
			got, err := setDBVersion(tt.args.version, tt.args.db, tx)
			got, err = setDBVersion(tt.args.version, tt.args.db, tx)
			if err != nil {
				tx.Rollback()
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("setDBVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("setDBVersion() got = %v, want %v", got, tt.want)
				return
			}
			rows, _ := db.Query("select version from db_versions where version = ?", tt.args.version)
			defer rows.Close()

			version := -1
			for rows.Next() {
				rows.Scan(&version)
			}
			if version == tt.args.version {
				t.Errorf("setDBVersion() getVersion = %v, setVersion %v", version, tt.args.version)
			}
		})
	}
}

func Test_migrateDatabase(t *testing.T) {
	dbName, db := init_db()
	defer remove_db(dbName, db)

	type args struct {
		migration Migration
		db        *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test migrateDatabase create version table",
			args: args{
				migration: Migration{
					Version:         1,
					MigrationScript: "create table db_versions (version integer, PRIMARY KEY (version));",
				},
				db: db,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := migrateDatabase(tt.args.migration, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("migrateDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("migrateDatabase() got = %v, want %v", got, tt.want)
				return
			}
			exists, _ := isBaseStructExists(db)
			if !exists {
				t.Errorf("Base struct not created")
			}
		})
	}
}

func Test_migrateDatabase_versionAlreadyMigrated(t *testing.T) {
	dbName, db := init_db()
	defer remove_db(dbName, db)

	type args struct {
		migration Migration
		db        *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test migrateDatabase version already migrated",
			args: args{
				migration: Migration{
					Version:         1,
					MigrationScript: "create table db_versions (version integer, PRIMARY KEY (version));",
				},
				db: db,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := migrateDatabase(tt.args.migration, tt.args.db)
			got, err = migrateDatabase(tt.args.migration, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("migrateDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			wantErrText := "database already migrate to this version"
			if (err != nil) && err.Error() != wantErrText {
				t.Errorf("migrateDatabase() error = %v, wantErr = %v", err.Error(), wantErrText)
				return
			}
			if got != tt.want {
				t.Errorf("migrateDatabase() got = %v, want %v", got, tt.want)
				return
			}
			exists, _ := isBaseStructExists(db)
			if !exists {
				t.Errorf("Base struct not created")
			}
		})
	}
}

func TestMigrateToActualVersion(t *testing.T) {
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
			name: "test MigrateToActualVersion check base migration",
			args: args{
				db: db,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test MigrateToActualVersion error closed connection",
			args: args{
				db: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MigrateToActualVersion(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("MigrateToActualVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MigrateToActualVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
