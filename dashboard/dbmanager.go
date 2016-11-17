package dashboard

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseManager interface {
	GetAll(pageSize, offset int) (result []Widget, err error)
	GetPages() (result []Page, err error)
	InsertWidget(widget *Widget) (int64, error)
	InsertOrUpdateWidget(widget *Widget) (int64, error)
	InsertWidgetToPage(pageId int64, widgetId string) (int64, error)
	InsertPage(page *Page) (int64, error)
	UpdatePage(page *Page) (int64, error)
	GetPageWidgets(pageId int64) (result []Widget, err error)
	GetUnlinkedWidgets() (result []Widget, err error)
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
	res, err := db.Exec(insertQuery, widget.Id, widget.Url, widget.Width, widget.Height)
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
	res, err := db.Exec(insertOrReplace, widget.Id, widget.Url, widget.Width, widget.Height)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (m *DbManager) GetAll(pageSize, offset int) (result []Widget, err error) {
	result = []Widget{}
	selectAllWithPaging := `
		SELECT * FROM widgets LIMIT ?,?;
	`
	db, err := m.Db()
	if err != nil {
		return
	}
	rows, err := db.Query(selectAllWithPaging, offset, pageSize)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var row Widget
		err = rows.Scan(&row.Id, &row.Width, &row.Height, &row.Url, &row.Content)
		if err != nil {
			return
		}
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

func (m *DbManager) InsertWidgetToPage(pageId int64, widgetId string) (int64, error) {
	insertOrReplace := `
		INSERT INTO page_widgets (id_widget,id_page)
		VALUES (?,?);
	`
	db, err := m.Db()
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(insertOrReplace, widgetId, pageId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (m *DbManager) GetPageWidgets(pageId int64) (result []Widget, err error) {
	result = []Widget{}
	selectAllPageWidgets := `
		SELECT * FROM widgets WHERE id in (SELECT id_widget FROM page_widgets WHERE id_page=?);
	`
	db, err := m.Db()
	if err != nil {
		return
	}
	rows, err := db.Query(selectAllPageWidgets, pageId)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var row Widget
		err = rows.Scan(&row.Id, &row.Width, &row.Height, &row.Url, &row.Content)
		if err != nil {
			return
		}
		result = append(result, row)
	}
	return
}

func (m *DbManager) GetUnlinkedWidgets() (result []Widget, err error) {
	result = []Widget{}
	selectQuery := `
		SELECT * FROM widgets WHERE id not in (SELECT DISTINCT id_widget FROM page_widgets);
	`
	db, err := m.Db()
	if err != nil {
		return
	}
	rows, err := db.Query(selectQuery)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var row Widget
		err = rows.Scan(&row.Id, &row.Width, &row.Height, &row.Url, &row.Content)
		if err != nil {
			return
		}
		if row.Id != nil {
			result = append(result, row)
		}
	}
	return
}

func (m *DbManager) GetPages() (result []Page, err error) {
	result = []Page{}
	selectQuery := `
		SELECT p.*, w.* FROM pages as p
		LEFT OUTER JOIN page_widgets as pw ON p.id = pw.id_page
		LEFT OUTER JOIN widgets as w ON w.id = pw.id_widget
		ORDER BY p.id;
	`
	db, err := m.Db()
	if err != nil {
		return
	}
	rows, err := db.Query(selectQuery)
	if err != nil {
		return
	}
	defer rows.Close()
	var page *Page
	var currentPageId int64

	for rows.Next() {
		var rowWidget Widget
		var rowPage Page
		err = rows.Scan(&rowPage.Id,
			&rowPage.Title,
			&rowPage.Visible,
			&rowWidget.Id,
			&rowWidget.Width,
			&rowWidget.Height,
			&rowWidget.Url,
			&rowWidget.Content)

		if err != nil {
			return
		}

		if currentPageId != rowPage.Id {
			if page != nil {
				result = append(result, *page)
				page = nil
			}
			page = &Page{}
			page.Id = rowPage.Id
			page.Title = rowPage.Title
			page.Visible = rowPage.Visible
		}

		if rowWidget.Id != nil {
			page.Widgets = append(page.Widgets, rowWidget)
		}

		currentPageId = rowPage.Id
	}

	if page != nil {
		result = append(result, *page)
	}

	return
}
