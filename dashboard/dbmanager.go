package dashboard

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseManager interface {
	InsertWidget(widget *Widget) (int64, error)
	Close() error
}

func NewDbManager(path string) DatabaseManager {
	return &DbManager{path: path}
}

type DbManager struct {
	db *sql.DB
	path string
}

func (m *DbManager) Db() (*sql.DB, error) {
	var err error
	if m.db == nil {
		m.db, err = sql.Open("sqlite3", m.path)
		if err != nil {
			return nil, err
		}
	}
	return m.db, err
}

func (m *DbManager) Close() error {
	if m.db != nil {
		err := m.db.Close()
		return err
	}
	return nil
}

func (m *DbManager) InsertWidget(widget *Widget) (int64, error) {
	insertQuery := `
		INSERT INTO widgets (id,url,width,height)
		VALUES (?,?,?,?);
	`
	db, err := m.Db()
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(insertQuery, widget.ID, widget.Url, widget.Width, widget.Height)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}