package dashboard

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var (
	dbMigrator = NewMigrator("test.db")
	dbManager = NewDbManager("test.db")
)

func TestDbManager(t *testing.T) {
	errors, ok := dbMigrator.Up()
	if !ok {
		t.Fatal("Migrations Up: ", errors)
	}
	t.Run("Successfully insert widget info", func(t *testing.T) {
		var widget = &Widget{Url:"http://example.com", ID:"test_widget_1", Height:350, Width:400}
		count, err := dbManager.InsertWidget(widget)
		if err != nil {
			t.Fatal("ERROR:", err)
		}
		assert.Equal(t, int64(1), count)
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