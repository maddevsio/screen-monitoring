package dashboard

import (
	"fmt"

	_ "github.com/mattes/migrate/driver/sqlite3"
	"github.com/mattes/migrate/migrate"
)

type Migrator interface {
	Up() ([]error, bool)
	Down() ([]error, bool)
}

type DbMigration struct {
	dbname      string
	migratePath string
}

func (dbm *DbMigration) Up() ([]error, bool) {
	dbConnString := fmt.Sprintf("sqlite3://%s", dbm.dbname)
	return migrate.UpSync(dbConnString, "./migrations")
}

func (dbm *DbMigration) Down() ([]error, bool) {
	dbConnString := fmt.Sprintf("sqlite3://%s", dbm.dbname)
	return migrate.DownSync(dbConnString, "./migrations")
}

func NewMigrator(dbname, migratePath string) Migrator {
	return &DbMigration{
		dbname:      dbname,
		migratePath: migratePath,
	}
}
