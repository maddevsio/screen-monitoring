package models

import (
	"upper.io/db.v2/lib/sqlbuilder"
	"upper.io/db.v2/sqlite"
)

type Datastore interface {
	CountersCreate(*Counter) error
	LastCounters() (*Counter, error)
	CountersForLastMonth() ([]*Counter, error)
}

type DB struct {
	sqlbuilder.Database
}

func NewDB(dbPath string) (*DB, error) {
	settings := sqlite.ConnectionURL{
		Database: dbPath,
	}
	conn, err := sqlite.Open(settings)
	if err != nil {
		return nil, err
	}
	return &DB{conn}, nil
}
