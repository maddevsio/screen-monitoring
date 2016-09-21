package dashboard

import (
	"github.com/mattes/migrate/migrate"
	_ "github.com/mattes/migrate/driver/sqlite3"
	"fmt"
)


type Migrator interface {
	Up() ([]error, bool)
	Down() ([]error, bool)
}

type DbMigration struct {
	dbname string
}

func (dbm *DbMigration) Up() ([]error, bool) {
	dbConnString := fmt.Sprintf("sqlite3://%s", dbm.dbname)
	return migrate.UpSync(dbConnString, "./migrations")
}

func (dbm *DbMigration) Down() ([]error, bool) {
	dbConnString := fmt.Sprintf("sqlite3://%s", dbm.dbname)
	return migrate.DownSync(dbConnString, "./migrations")
}

func NewMigrator(dbname string) Migrator {
	return &DbMigration{dbname:dbname}
}