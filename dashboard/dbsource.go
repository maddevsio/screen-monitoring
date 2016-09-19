package dashboard

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DbManagerInstance *DbManager

func init(){
	DbManagerInstance = &DbManager{}
	DbManagerInstance.initDb()
}

type DbManager struct {
	db *sql.DB
}

func (m *DbManager) Db() (*sql.DB, error) {
	var err error
	if m.db == nil {
		m.db, err = sql.Open("sqlite3", "./screen_monitoring.db")
		if err != nil {
			return nil, err
		}
	}
	return m.db, err
}

func (m *DbManager) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *DbManager) initDb() error {
	createWidgetsTable := `
		CREATE TABLE IF NOT EXISTS widgets (
			id	TEXT NOT NULL UNIQUE,
			width	INTEGER NOT NULL DEFAULT 300,
			height	INTEGER NOT NULL DEFAULT 300,
			url	TEXT NOT NULL,
			content	TEXT,
			PRIMARY KEY(id)
		);
	`
	db, err := m.Db()
	if err != nil {
		return err
	}
	_, err = db.Exec(createWidgetsTable)
	return err
}

func (m *DbManager) InsertWidget(widget *Widget) (int64, error) {
	insertQuery := `
		INSERT INTO widgets (id,url,width,height)
		VALUES (?,?,?);
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