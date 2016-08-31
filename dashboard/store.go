package dashboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"sync"
	"time"
)

type WidgetStore struct {
	db *bolt.DB
	sync.RWMutex
}

func (w *WidgetStore) getDb() (*bolt.DB, error) {
	w.Lock()
	defer w.Unlock()
	var err error
	if w.db == nil {
		w.db, err = bolt.Open("dashboard.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			return nil, err
		}
		return w.db, nil
	}
	return w.db, nil
}

func (w *WidgetStore) Save(bucket string, widget *Widget) error {

	if w.db != nil {
		return w.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucket([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			jsonBytes, err := widget.ToJson()
			if err != nil {
				return err
			}
			return b.Put([]byte(widget.ID), jsonBytes)
		})
	}
	return errors.New("Database not initialized!")
}

func (w *WidgetStore) GetAll(bucket string) ([]Widget, error) {
	var widgets []Widget
	err := w.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.ForEach(func(k, v []byte) error {
			var widget Widget
			errJson := json.Unmarshal(v, &widget)
			if errJson == nil {
				widgets = append(widgets, widget)
			}
			return errJson
		})
	})
	return widgets, err
}

func (w *WidgetStore) Shutdown() error {
	if w.db != nil {
		return w.db.Close()
	}
	return errors.New("No Database to close")
}
