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
	t.Run("Successfully insert widget info", func(t *testing.T) {
		up(t)
		var widget = &Widget{Url: "http://example.com", Id: "test_widget_1", Height: 350, Width: 400}
		count, err := dbManager.InsertWidget(widget)
		if err != nil {
			t.Fatal("ERROR:", err)
		}
		assert.Equal(t, int64(1), count)
		teardown(t)
	})
	t.Run("Fail to insert widget with same name", func(t *testing.T) {
		up(t)
		var widget = &Widget{Url: "http://example.com", Id: "test_widget_1", Height: 350, Width: 400}
		var widget2 = &Widget{Url: "http://some.site.com", Id: "test_widget_1", Height: 450, Width: 350}
		dbManager.InsertWidget(widget)
		count, err := dbManager.InsertWidget(widget2)
		assert.NotNil(t, err)
		assert.Equal(t, "UNIQUE constraint failed: widgets.id", err.Error())
		assert.Equal(t, int64(0), count)
		teardown(t)
	})
	t.Run("Should update widget data if widget exists", func(t *testing.T) {
		up(t)
		var widget = &Widget{Url: "http://example.com", Id: "test_widget_1", Height: 350, Width: 400}
		var widget2 = &Widget{Url: "http://some.site.com", Id: "test_widget_1", Height: 450, Width: 310}
		dbManager.InsertWidget(widget)
		id, err := dbManager.InsertOrUpdateWidget(widget2)
		assert.Nil(t, err)
		assert.True(t, id > 0)
		teardown(t)
	})
	t.Run("Should return all widgets", func(t *testing.T) {
		up(t)
		var widget = Widget{Url: "http://example.com", Id: "test_widget_1", Height: 350, Width: 400}
		var widget2 = Widget{Url: "http://some.site.com", Id: "test_widget_2", Height: 450, Width: 310}
		_, err := dbManager.InsertWidget(&widget)
		assert.Nil(t, err)
		_, err = dbManager.InsertWidget(&widget2)
		assert.Nil(t, err)
		actual, err := dbManager.GetAll(10, 0)
		assert.Nil(t, err)
		assert.Equal(t, widget, actual[0])
		assert.Equal(t, widget2, actual[1])
		teardown(t)
	})
	t.Run("Should return valid count rows per page", func(t *testing.T) {
		up(t)
		var widget1 = Widget{Url: "http://example.com", Id: "l_test_widget_1", Height: 350, Width: 400}
		var widget2 = Widget{Url: "http://some.site1.com", Id: "j_test_widget_2", Height: 450, Width: 310}
		var widget3 = Widget{Url: "http://some.site2.com", Id: "a_test_widget_3", Height: 450, Width: 310}
		var widget4 = Widget{Url: "http://some.site3.com", Id: "s_test_widget_4", Height: 450, Width: 310}
		var widget5 = Widget{Url: "http://some.site4.com", Id: "f_test_widget_5", Height: 450, Width: 310}
		var widget6 = Widget{Url: "http://some.site5.com", Id: "e_test_widget_6", Height: 450, Width: 310}
		var widget7 = Widget{Url: "http://some.site6.com", Id: "b_test_widget_7", Height: 450, Width: 310}
		var widget8 = Widget{Url: "http://some.site7.com", Id: "c_test_widget_8", Height: 450, Width: 310}
		var all = []Widget{widget1, widget2, widget3, widget4, widget5, widget6, widget7, widget8}
		var expected = []Widget{widget1, widget2, widget3, widget4}
		for _, widget := range all {
			dbManager.InsertWidget(&widget)
		}
		actual, err := dbManager.GetAll(4, 0)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
		teardown(t)
	})

	t.Run("Should return empty array if no wigets in db", func(t *testing.T) {
		up(t)
		actual, err := dbManager.GetAll(4, 0)
		assert.Nil(t, err)
		assert.Empty(t, actual)
		teardown(t)
	})

	t.Run("Should success create page", func(t *testing.T) {
		up(t)
		var page = &Page{Title: "Page 1", Visible: true}
		count, err := dbManager.InsertPage(page)
		assert.Equal(t, int64(1), count)
		assert.Nil(t, err)
		teardown(t)
	})

	t.Run("Should not insert page with same title", func(t *testing.T) {
		up(t)
		var page = &Page{Title: "Page 1", Visible: true}
		var page2 = &Page{Title: "Page 1", Visible: true}
		dbManager.InsertPage(page)
		count, err := dbManager.InsertPage(page2)
		assert.NotNil(t, err)
		assert.Equal(t, "UNIQUE constraint failed: pages.title", err.Error())
		assert.Equal(t, int64(0), count)
		teardown(t)
	})

	t.Run("Should update page by id", func(t *testing.T) {
		up(t)
		var page = &Page{Title: "Page 2", Visible: true}
		id, err := dbManager.InsertPage(page)
		page.Id = id
		page.Visible = false
		page.Title = "Page title changed"
		count, err := dbManager.UpdatePage(page)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
		teardown(t)
	})

	t.Run("Should link widget to page", func(t *testing.T) {
		up(t)
		var page = &Page{Title: "Page 3", Visible: true}
		var widget = &Widget{Url: "http://example1.com", Id: "widget_page_3", Height: 450, Width: 300}
		pid, err := dbManager.InsertPage(page)
		if err != nil {
			t.Fatal("Error creating page: ", err)
		}
		_, err = dbManager.InsertWidget(widget)
		if err != nil {
			t.Fatal("Error creating widget: ", err)
		}

		count, err := dbManager.InsertWidgetToPage(pid, widget.Id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
		teardown(t)
	})

	t.Run("Should link widget to page twice with same widget id", func(t *testing.T) {
		up(t)
		var page = &Page{Title: "Page 3", Visible: true}
		var widget = &Widget{Url: "http://example1.com", Id: "widget_page_3", Height: 450, Width: 300}
		pid, err := dbManager.InsertPage(page)
		if err != nil {
			t.Fatal("Error creating page: ", err)
		}
		_, err = dbManager.InsertWidget(widget)
		if err != nil {
			t.Fatal("Error creating widget: ", err)
		}

		_, err = dbManager.InsertWidgetToPage(pid, widget.Id)
		count, err := dbManager.InsertWidgetToPage(pid, widget.Id)
		assert.NotNil(t, err)
		assert.Equal(t, PAGE_WIDGET_UNIQUE_ERROR, err.Error())
		assert.Equal(t, int64(0), count)
		teardown(t)
	})

	t.Run("Get widgets by page Id", func(t *testing.T) {
		up(t)
		var page = &Page{Title: "Page 1", Visible: true}
		var page2 = &Page{Title: "Page 2", Visible: true}
		var widget = Widget{Url: "http://example1.com", Id: "widget_page_1", Height: 450, Width: 300}
		var widget2 = Widget{Url: "http://example2.com", Id: "widget_page_2", Height: 420, Width: 200}
		var widget3 = Widget{Url: "http://example3.com", Id: "widget_page_3", Height: 320, Width: 210}
		pid, err := dbManager.InsertPage(page)
		pid2, err := dbManager.InsertPage(page2)
		dbManager.InsertWidget(&widget)
		dbManager.InsertWidget(&widget2)
		dbManager.InsertWidget(&widget3)

		_, err = dbManager.InsertWidgetToPage(pid, widget.Id)
		t.Log(err)
		_, err = dbManager.InsertWidgetToPage(pid, widget2.Id)
		t.Log(err)
		_, err = dbManager.InsertWidgetToPage(pid2, widget3.Id)
		t.Log(err)

		expected := []Widget{widget, widget2}
		expected2 := []Widget{widget3}

		actual, err := dbManager.GetPageWidgets(pid)
		actual2, err := dbManager.GetPageWidgets(pid2)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
		assert.Equal(t, expected2, actual2)
		teardown(t)
	})

	t.Run("Should return unlinked widgets to pages", func(t *testing.T) {
		up(t)
		var page = Page{Title: "Page 1", Visible: true}
		var widget = Widget{Url: "http://example1.com", Id: "widget_page_1", Height: 450, Width: 300}
		var widget2 = Widget{Url: "http://example2.com", Id: "widget_page_2", Height: 420, Width: 200}
		pid, err := dbManager.InsertPage(&page)
		dbManager.InsertWidget(&widget)
		dbManager.InsertWidget(&widget2)
		expected := []Widget{widget2}
		_, err = dbManager.InsertWidgetToPage(pid, widget.Id)
		actual, err := dbManager.GetUnlinkedWidgets()
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
		teardown(t)
	})

	t.Run("Should return empty array if no widgets unlinked", func(t *testing.T) {
		up(t)
		expected := []Widget{}
		actual, err := dbManager.GetUnlinkedWidgets()
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
		teardown(t)
	})

	t.Run("Should return empty array if no unlinked widgets to pages", func(t *testing.T) {
		up(t)
		var page = Page{Title: "Page 1", Visible: true}
		var widget = Widget{Url: "http://example1.com", Id: "widget_page_1", Height: 450, Width: 300}
		var widget2 = Widget{Url: "http://example2.com", Id: "widget_page_2", Height: 420, Width: 200}
		pid, err := dbManager.InsertPage(&page)
		dbManager.InsertWidget(&widget)
		dbManager.InsertWidget(&widget2)
		expected := []Widget{}
		_, err = dbManager.InsertWidgetToPage(pid, widget.Id)
		_, err = dbManager.InsertWidgetToPage(pid, widget2.Id)
		actual, err := dbManager.GetUnlinkedWidgets()
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
		teardown(t)
	})

	t.Run("Should return pages with widgets", func(t *testing.T) {
		up(t)
		var page = Page{Title: "Page 1", Visible: true}
		var widget = Widget{Url: "http://example3.com", Id: "widget_page_3", Height: 450, Width: 300}
		var widget2 = Widget{Url: "http://example4.com", Id: "widget_page_4", Height: 420, Width: 200}
		pid, err := dbManager.InsertPage(&page)
		dbManager.InsertWidget(&widget)
		dbManager.InsertWidget(&widget2)
		expected := []Page{Page{
			Id:      pid,
			Title:   page.Title,
			Visible: page.Visible,
			Widgets: []Widget{widget, widget2},
		}}
		_, err = dbManager.InsertWidgetToPage(pid, widget.Id)
		_, err = dbManager.InsertWidgetToPage(pid, widget2.Id)

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
