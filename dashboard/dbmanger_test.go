package dashboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dbMigrator = NewMigrator("test.db", "./migrations")
	dbManager  = NewDbManager("test.db")
)

const (
	PAGE_WIDGET_UNIQUE_ERROR = "UNIQUE constraint failed: page_widgets.id_widget, page_widgets.id_page"
)

func teardown(t *testing.T) {
	errors, ok := dbMigrator.Down()
	t.Log("Down migrations")
	if !ok {
		t.Fatal("Migrations Down: ", errors)
	}
}

func up(t *testing.T) {
	errors, ok := dbMigrator.Up()
	t.Log("Up migrations")
	if !ok {
		t.Fatal("Migrations Up: ", errors)
	}
}

func TestDbManager(t *testing.T) {
	// t.Run("Successfully insert widget info", func(t *testing.T) {
	// 	up(t)
	// 	var widget = &Widget{
	// 		Url:     sql.NullString{String: "http://example.com", Valid: true},
	// 		Id:      sql.NullString{String: "test_widget_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 350, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 400, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	count, err := dbManager.InsertWidget(widget)
	// 	if err != nil {
	// 		t.Fatal("ERROR:", err)
	// 	}
	// 	assert.Equal(t, int64(1), count)
	// 	teardown(t)
	// })
	// t.Run("Fail to insert widget with same name", func(t *testing.T) {
	// 	up(t)
	// 	var widget = &Widget{
	// 		Url:     sql.NullString{String: "http://example.com", Valid: true},
	// 		Id:      sql.NullString{String: "test_widget_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 350, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 400, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget2 = &Widget{
	// 		Url:     sql.NullString{String: "http://some.site.com", Valid: true},
	// 		Id:      sql.NullString{String: "test_widget_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 350, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	dbManager.InsertWidget(widget)
	// 	count, err := dbManager.InsertWidget(widget2)
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, "UNIQUE constraint failed: widgets.id", err.Error())
	// 	assert.Equal(t, int64(0), count)
	// 	teardown(t)
	// })
	// t.Run("Fail to insert widget with negative coordinates", func(t *testing.T) {
	// 	up(t)
	// 	var widget = &Widget{
	// 		Url:     sql.NullString{String: "http://example.com", Valid: true},
	// 		Id:      sql.NullString{String: "test_widget_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 350, Valid: true},
	// 		Width:   sql.NullInt64{Int64: -400, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	count, err := dbManager.InsertWidget(widget)
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, "CHECK constraint failed: widgets", err.Error())
	// 	assert.Equal(t, int64(0), count)
	// 	teardown(t)
	// })
	// t.Run("Fail to insert widget with empty url", func(t *testing.T) {
	// 	up(t)
	// 	var widget = &Widget{
	// 		Url:     sql.NullString{String: "", Valid: true},
	// 		Id:      sql.NullString{String: "test_widget_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 350, Valid: true},
	// 		Width:   sql.NullInt64{Int64: -400, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	count, err := dbManager.InsertWidget(widget)
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, "CHECK constraint failed: widgets", err.Error())
	// 	assert.Equal(t, int64(0), count)
	// 	teardown(t)
	// })
	// t.Run("Should update widget data if widget exists", func(t *testing.T) {
	// 	up(t)
	// 	var widget = &Widget{
	// 		Url:     sql.NullString{String: "http://example.com", Valid: true},
	// 		Id:      sql.NullString{String: "test_widget_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 350, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 400, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget2 = &Widget{
	// 		Url:     sql.NullString{String: "http://some.site.com", Valid: true},
	// 		Id:      sql.NullString{String: "test_widget_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 310, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	dbManager.InsertWidget(widget)
	// 	id, err := dbManager.InsertOrUpdateWidget(widget2)
	// 	assert.Nil(t, err)
	// 	assert.True(t, id > 0)
	// 	teardown(t)
	// })
	// t.Run("Should return all widgets", func(t *testing.T) {
	// 	up(t)
	// 	var widget = Widget{
	// 		Url:     sql.NullString{String: "http://example.com", Valid: true},
	// 		Id:      sql.NullString{String: "test_widget_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 350, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 400, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget2 = Widget{
	// 		Url:     sql.NullString{String: "http://some.site.com", Valid: true},
	// 		Id:      sql.NullString{String: "test_widget_2", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 310, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	_, err := dbManager.InsertWidget(&widget)
	// 	assert.Nil(t, err)
	// 	_, err = dbManager.InsertWidget(&widget2)
	// 	assert.Nil(t, err)
	// 	actual, err := dbManager.GetAll(10, 0)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, widget, actual[0])
	// 	assert.Equal(t, widget2, actual[1])
	// 	teardown(t)
	// })
	// t.Run("Should return valid count rows per page", func(t *testing.T) {
	// 	up(t)
	// 	var widget1 = Widget{
	// 		Url:     sql.NullString{String: "http://example.com", Valid: true},
	// 		Id:      sql.NullString{String: "l_test_widget_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 350, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 400, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget2 = Widget{
	// 		Url:     sql.NullString{String: "http://some.site1.com", Valid: true},
	// 		Id:      sql.NullString{String: "j_test_widget_2", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 310, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget3 = Widget{
	// 		Url:     sql.NullString{String: "http://some.site2.com", Valid: true},
	// 		Id:      sql.NullString{String: "a_test_widget_3", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 310, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget4 = Widget{
	// 		Url:     sql.NullString{String: "http://some.site3.com", Valid: true},
	// 		Id:      sql.NullString{String: "s_test_widget_4", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 310, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget5 = Widget{
	// 		Url:     sql.NullString{String: "http://some.site4.com", Valid: true},
	// 		Id:      sql.NullString{String: "f_test_widget_5", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 310, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget6 = Widget{
	// 		Url:     sql.NullString{String: "http://some.site5.com", Valid: true},
	// 		Id:      sql.NullString{String: "e_test_widget_6", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 310, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget7 = Widget{
	// 		Url:     sql.NullString{String: "http://some.site6.com", Valid: true},
	// 		Id:      sql.NullString{String: "b_test_widget_7", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 310, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget8 = Widget{
	// 		Url:     sql.NullString{String: "http://some.site7.com", Valid: true},
	// 		Id:      sql.NullString{String: "c_test_widget_8", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 310, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var all = []Widget{widget1, widget2, widget3, widget4, widget5, widget6, widget7, widget8}
	// 	var expected = []Widget{widget1, widget2, widget3, widget4}
	// 	for _, widget := range all {
	// 		dbManager.InsertWidget(&widget)
	// 	}
	// 	actual, err := dbManager.GetAll(4, 0)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, expected, actual)
	// 	teardown(t)
	// })
	//
	// t.Run("Should return empty array if no wigets in db", func(t *testing.T) {
	// 	up(t)
	// 	actual, err := dbManager.GetAll(4, 0)
	// 	assert.Nil(t, err)
	// 	assert.Empty(t, actual)
	// 	teardown(t)
	// })
	//
	// t.Run("Should success create page", func(t *testing.T) {
	// 	up(t)
	// 	var page = &Page{Title: "Page 1", Visible: true}
	// 	count, err := dbManager.InsertPage(page)
	// 	assert.Equal(t, int64(1), count)
	// 	assert.Nil(t, err)
	// 	teardown(t)
	// })
	//
	// t.Run("Should not create page with empty title", func(t *testing.T) {
	// 	up(t)
	// 	var page = &Page{Title: "", Visible: true}
	// 	count, err := dbManager.InsertPage(page)
	// 	assert.Equal(t, int64(0), count)
	// 	assert.NotNil(t, err)
	// 	teardown(t)
	// })
	//
	// t.Run("Should not insert page with same title", func(t *testing.T) {
	// 	up(t)
	// 	var page = &Page{Title: "Page 1", Visible: true}
	// 	var page2 = &Page{Title: "Page 1", Visible: true}
	// 	dbManager.InsertPage(page)
	// 	count, err := dbManager.InsertPage(page2)
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, "UNIQUE constraint failed: pages.title", err.Error())
	// 	assert.Equal(t, int64(0), count)
	// 	teardown(t)
	// })
	//
	// t.Run("Should update page by id", func(t *testing.T) {
	// 	up(t)
	// 	var page = &Page{Title: "Page 2", Visible: true}
	// 	id, err := dbManager.InsertPage(page)
	// 	page.Id = id
	// 	page.Visible = false
	// 	page.Title = "Page title changed"
	// 	count, err := dbManager.UpdatePage(page)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, int64(1), count)
	// 	teardown(t)
	// })
	//
	// t.Run("Should link widget to page", func(t *testing.T) {
	// 	up(t)
	// 	var page = &Page{Title: "Page 3", Visible: true}
	// 	var widget = &Widget{
	// 		Url:     sql.NullString{String: "http://example1.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_3", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 300, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	pid, err := dbManager.InsertPage(page)
	// 	if err != nil {
	// 		t.Fatal("Error creating page: ", err)
	// 	}
	// 	_, err = dbManager.InsertWidget(widget)
	// 	if err != nil {
	// 		t.Fatal("Error creating widget: ", err)
	// 	}
	//
	// 	count, err := dbManager.InsertWidgetToPage(pid, widget.Id.String)
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, int64(1), count)
	// 	teardown(t)
	// })
	//
	// t.Run("Should link widget to page twice with same widget id", func(t *testing.T) {
	// 	up(t)
	// 	var page = &Page{Title: "Page 3", Visible: true}
	// 	var widget = &Widget{
	// 		Url:     sql.NullString{String: "http://example1.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_3", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 300, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	pid, err := dbManager.InsertPage(page)
	// 	if err != nil {
	// 		t.Fatal("Error creating page: ", err)
	// 	}
	// 	_, err = dbManager.InsertWidget(widget)
	// 	if err != nil {
	// 		t.Fatal("Error creating widget: ", err)
	// 	}
	//
	// 	_, err = dbManager.InsertWidgetToPage(pid, widget.Id.String)
	// 	count, err := dbManager.InsertWidgetToPage(pid, widget.Id.String)
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, PAGE_WIDGET_UNIQUE_ERROR, err.Error())
	// 	assert.Equal(t, int64(0), count)
	// 	teardown(t)
	// })
	//
	// t.Run("Get widgets by page Id", func(t *testing.T) {
	// 	up(t)
	// 	var page = &Page{Title: "Page 1", Visible: true}
	// 	var page2 = &Page{Title: "Page 2", Visible: true}
	// 	var widget = Widget{
	// 		Url:     sql.NullString{String: "http://example1.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 300, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget2 = Widget{
	// 		Url:     sql.NullString{String: "http://example2.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_2", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 420, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 200, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget3 = Widget{
	// 		Url:     sql.NullString{String: "http://example3.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_3", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 320, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 210, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	pid, err := dbManager.InsertPage(page)
	// 	pid2, err := dbManager.InsertPage(page2)
	// 	dbManager.InsertWidget(&widget)
	// 	dbManager.InsertWidget(&widget2)
	// 	dbManager.InsertWidget(&widget3)
	//
	// 	_, err = dbManager.InsertWidgetToPage(pid, widget.Id.String)
	// 	t.Log(err)
	// 	_, err = dbManager.InsertWidgetToPage(pid, widget2.Id.String)
	// 	t.Log(err)
	// 	_, err = dbManager.InsertWidgetToPage(pid2, widget3.Id.String)
	// 	t.Log(err)
	//
	// 	expected := []Widget{widget, widget2}
	// 	expected2 := []Widget{widget3}
	//
	// 	actual, err := dbManager.GetPageWidgets(pid)
	// 	actual2, err := dbManager.GetPageWidgets(pid2)
	//
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, expected, actual)
	// 	assert.Equal(t, expected2, actual2)
	// 	teardown(t)
	// })
	//
	// t.Run("Should return unlinked widgets to pages", func(t *testing.T) {
	// 	up(t)
	// 	var page = Page{Title: "Page 1", Visible: true}
	// 	var widget = Widget{
	// 		Url:     sql.NullString{String: "http://example1.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 300, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget2 = Widget{
	// 		Url:     sql.NullString{String: "http://example2.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_2", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 420, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 200, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	pid, err := dbManager.InsertPage(&page)
	// 	dbManager.InsertWidget(&widget)
	// 	dbManager.InsertWidget(&widget2)
	// 	expected := []Widget{widget2}
	// 	_, err = dbManager.InsertWidgetToPage(pid, widget.Id.String)
	// 	actual, err := dbManager.GetUnlinkedWidgets()
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, expected, actual)
	// 	teardown(t)
	// })
	//
	// t.Run("Should return empty array if no widgets unlinked", func(t *testing.T) {
	// 	up(t)
	// 	expected := []Widget{}
	// 	actual, err := dbManager.GetUnlinkedWidgets()
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, expected, actual)
	// 	teardown(t)
	// })
	//
	// t.Run("Should return empty array if no unlinked widgets to pages", func(t *testing.T) {
	// 	up(t)
	// 	var page = Page{Title: "Page 1", Visible: true}
	// 	var widget = Widget{
	// 		Url:     sql.NullString{String: "http://example1.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_1", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 300, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget2 = Widget{
	// 		Url:     sql.NullString{String: "http://example2.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_2", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 420, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 200, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	pid, err := dbManager.InsertPage(&page)
	// 	dbManager.InsertWidget(&widget)
	// 	dbManager.InsertWidget(&widget2)
	// 	expected := []Widget{}
	// 	_, err = dbManager.InsertWidgetToPage(pid, widget.Id.String)
	// 	_, err = dbManager.InsertWidgetToPage(pid, widget2.Id.String)
	// 	actual, err := dbManager.GetUnlinkedWidgets()
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, expected, actual)
	// 	teardown(t)
	// })
	//
	// t.Run("Should return pages with widgets", func(t *testing.T) {
	// 	up(t)
	// 	var page = Page{Title: "Page 1", Visible: true}
	// 	var widget = Widget{
	// 		Url:     sql.NullString{String: "http://example3.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_3", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 450, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 300, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	var widget2 = Widget{
	// 		Url:     sql.NullString{String: "http://example4.com", Valid: true},
	// 		Id:      sql.NullString{String: "widget_page_4", Valid: true},
	// 		Height:  sql.NullInt64{Int64: 420, Valid: true},
	// 		Width:   sql.NullInt64{Int64: 200, Valid: true},
	// 		Content: sql.NullString{String: "", Valid: true},
	// 	}
	// 	pid, err := dbManager.InsertPage(&page)
	// 	dbManager.InsertWidget(&widget)
	// 	dbManager.InsertWidget(&widget2)
	// 	expected := []Page{Page{
	// 		Id:      pid,
	// 		Title:   page.Title,
	// 		Visible: page.Visible,
	// 		Widgets: []Widget{widget, widget2},
	// 	}}
	// 	_, err = dbManager.InsertWidgetToPage(pid, widget.Id.String)
	// 	_, err = dbManager.InsertWidgetToPage(pid, widget2.Id.String)
	//
	// 	actual, err := dbManager.GetPages()
	// 	assert.Nil(t, err)
	// 	assert.Equal(t, expected, actual)
	// 	teardown(t)
	// })

	t.Run("Should return more than one pages with widgets", func(t *testing.T) {
		up(t)
		var page = Page{Title: "Page 1", Visible: true}
		var page2 = Page{Title: "Page 2", Visible: true}
		var widget = NewWidget(
			"widget_page_3",
			300,
			450,
			"http://example3.com",
		)
		var widget2 = NewWidget(
			"widget_page_4",
			420,
			200,
			"http://example4.com",
		)
		var widget3 = NewWidget(
			"widget_page_5",
			410,
			190,
			"http://example5.com",
		)
		pid, err := dbManager.InsertPage(&page)
		pid2, err := dbManager.InsertPage(&page2)
		dbManager.InsertWidget(&widget)
		dbManager.InsertWidget(&widget2)
		dbManager.InsertWidget(&widget3)
		expected := []Page{
			Page{
				Id:      pid,
				Title:   page.Title,
				Visible: page.Visible,
				Widgets: []Widget{widget, widget2},
			},
			Page{
				Id:      pid2,
				Title:   page2.Title,
				Visible: page2.Visible,
				Widgets: []Widget{widget3},
			},
		}
		_, err = dbManager.InsertWidgetToPage(pid, *widget.Id)
		_, err = dbManager.InsertWidgetToPage(pid, *widget2.Id)
		_, err = dbManager.InsertWidgetToPage(pid2, *widget3.Id)

		actual, err := dbManager.GetPages()
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
		teardown(t)
	})

	err := dbManager.Close()
	if err != nil {
		t.Fatal("DB MANAGER: ", err)
	}
}
