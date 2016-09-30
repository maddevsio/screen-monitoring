package dashboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dbMigrator = NewMigrator("test.db", "./migrations")
	dbManager  = NewDbManager("test.db")
)

func TestDbManager(t *testing.T) {
	errors, ok := dbMigrator.Up()
	if !ok {
		t.Fatal("Migrations Up: ", errors)
	}
	t.Run("Successfully insert widget info", func(t *testing.T) {
		var widget = &Widget{Url: "http://example.com", ID: "test_widget_1", Height: 350, Width: 400}
		count, err := dbManager.InsertWidget(widget)
		if err != nil {
			t.Fatal("ERROR:", err)
		}
		assert.Equal(t, int64(1), count)
	})
	t.Run("Fail to insert widget with same name", func(t *testing.T) {
		var widget = &Widget{Url: "http://example.com", ID: "test_widget_1", Height: 350, Width: 400}
		var widget2 = &Widget{Url: "http://some.site.com", ID: "test_widget_1", Height: 450, Width: 350}
		dbManager.InsertWidget(widget)
		count, err := dbManager.InsertWidget(widget2)
		assert.NotNil(t, err)
		assert.Equal(t, "UNIQUE constraint failed: widgets.id", err.Error())
		assert.Equal(t, int64(0), count)
	})
	t.Run("Should update didget data if widget exists", func(t *testing.T) {
		var widget = &Widget{Url: "http://example.com", ID: "test_widget_1", Height: 350, Width: 400}
		var widget2 = &Widget{Url: "http://some.site.com", ID: "test_widget_1", Height: 450, Width: 310}
		dbManager.InsertWidget(widget)
		count, err := dbManager.InsertOrUpdateWidget(widget2)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	})
	t.Run("Should return all widgets", func(t *testing.T) {
		var widget = &Widget{Url: "http://example.com", ID: "test_widget_1", Height: 350, Width: 400}
		var widget2 = &Widget{Url: "http://some.site.com", ID: "test_widget_2", Height: 450, Width: 310}
		var expected = []*Widget{widget, widget2}
		dbManager.InsertWidget(widget)
		dbManager.InsertWidget(widget2)
		actual, err := dbManager.GetAll(10, 0)
		assert.Nil(t, err)
		assert.InEpsilonSlice(t, expected, actual, float64(0))
	})
	t.Run("Should return valid count rows per page", func(t *testing.T) {
		var widget1 = &Widget{Url: "http://example.com", ID: "l_test_widget_1", Height: 350, Width: 400}
		var widget2 = &Widget{Url: "http://some.site1.com", ID: "j_test_widget_2", Height: 450, Width: 310}
		var widget3 = &Widget{Url: "http://some.site2.com", ID: "a_test_widget_3", Height: 450, Width: 310}
		var widget4 = &Widget{Url: "http://some.site3.com", ID: "s_test_widget_4", Height: 450, Width: 310}
		var widget5 = &Widget{Url: "http://some.site4.com", ID: "f_test_widget_5", Height: 450, Width: 310}
		var widget6 = &Widget{Url: "http://some.site5.com", ID: "e_test_widget_6", Height: 450, Width: 310}
		var widget7 = &Widget{Url: "http://some.site6.com", ID: "b_test_widget_7", Height: 450, Width: 310}
		var widget8 = &Widget{Url: "http://some.site7.com", ID: "c_test_widget_8", Height: 450, Width: 310}
		var all = []*Widget{widget1, widget2, widget3, widget4, widget5, widget6, widget7, widget8}
		var expected = []*Widget{widget1, widget2, widget3, widget4}
		for _, widget := range all {
			dbManager.InsertWidget(widget)
		}
		actual, err := dbManager.GetAll(4, 0)
		assert.Nil(t, err)
		assert.InEpsilonSlice(t, expected, actual, float64(0))
	})

	t.Run("Should return empty array if now wigets in db", func(t *testing.T) {
		actual, err := dbManager.GetAll(4, 0)
		assert.Nil(t, err)
		assert.Empty(t, actual)
	})

	t.Run("Should success create page", func(t *testing.T) {
		var page = &Page{Title: "Page 1", Visible: true}
		count, err := dbManager.InsertPage(page)
		assert.Equal(t, int64(1), count)
		assert.Nil(t, err)
	})

	t.Run("Should not insert page with same title", func(t *testing.T) {
		var page = &Page{Title: "Page 1", Visible: true}
		var page2 = &Page{Title: "Page 1", Visible: true}
		dbManager.InsertPage(page)
		count, err := dbManager.InsertPage(page2)
		assert.NotNil(t, err)
		assert.Equal(t, "UNIQUE constraint failed: pages.title", err.Error())
		assert.Equal(t, int64(0), count)
	})

	errors, ok = dbMigrator.Down()
	if !ok {
		t.Fatal("Migrations Down: ", errors)
	}
	err := dbManager.Close()
	if err != nil {
		t.Fatal("DB MANAGER: ", err)
	}
}
