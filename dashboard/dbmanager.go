package dashboard

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseManager interface {
	GetAll(pageSize, offset int) (result []Widget, err error)
	InsertWidget(widget *Widget) (int64, error)
	InsertOrUpdateWidget(widget *Widget) (int64, error)
	InsertPage(page *Page) (int64, error)
	UpdatePage(page *Page) (int64, error)
	Close() error
}

func NewDbManager(path string) DatabaseManager {
	return &DbManager{path: path}
}

type DbManager struct {
	db   *sql.DB
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

func (m *DbManager) InsertOrUpdateWidget(widget *Widget) (int64, error) {
	insertOrReplace := `
		INSERT OR REPLACE INTO widgets (id,url,width,height)
		VALUES (?,?,?,?);
	`
	db, err := m.Db()
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(insertOrReplace, widget.ID, widget.Url, widget.Width, widget.Height)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (m *DbManager) GetAll(pageSize, offset int) (result []Widget, err error) {
	selectAllWithPaging := `
		SELECT * FROM widgets LIMIT ?,?;
	`
	db, err := m.Db()
	if err != nil {
		return
	}
	rows, err := db.Query(selectAllWithPaging, pageSize, offset)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var row Widget
		rows.Scan(&row)
		result = append(result, row)
	}
	return
}

func (m *DbManager) InsertPage(page *Page) (int64, error) {
	insertOrReplace := `
		INSERT INTO pages (title,visible)
		VALUES (?,?);
	`
	db, err := m.Db()
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(insertOrReplace, page.Title, page.Visible)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (m *DbManager) UpdatePage(page *Page) (int64, error) {
	insertOrReplace := `
		UPDATE pages SET title=?, visible=? WHERE id = ?
	`
	db, err := m.Db()
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(insertOrReplace, page.Title, page.Visible, page.Id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
