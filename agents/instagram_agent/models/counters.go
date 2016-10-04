package models

import (
	"log"
	"time"
)

type Counter struct {
	Created    time.Time `db:"created"`
	Username   string    `db:"username"`
	Media      int64     `db:"media"`
	Follows    int64     `db:"follows"`
	FollowedBy int64     `db:"followed_by"`
}

func (db *DB) CountersCreate(c *Counter) error {
	countersCollection := db.Collection("counters")
	counterID, err := countersCollection.Insert(Counter{
		Created:    time.Now(),
		Username:   c.Username,
		Media:      c.Media,
		Follows:    c.Follows,
		FollowedBy: c.FollowedBy,
	})

	if err != nil {
		log.Fatal("Inserting counters error: ", err)
		return err
	}
	log.Printf("Add new counters with ID: %d", counterID)
	return nil
}

func (db *DB) LastCounters() (*Counter, error) {
	var counter *Counter
	q := db.SelectFrom("counters").OrderBy("created DESC").Limit(1)
	err := q.One(&counter)
	if err != nil {
		return nil, err
	}
	return counter, nil
}

func (db *DB) CountersForLastMonth() ([]*Counter, error) {
	var counters []*Counter
	q := db.SelectFrom("counters").Where("created > datetime('created', '-1 month')")
	err := q.All(&counters)
	if err != nil {
		log.Fatal("Err getCounterForLastMonth", err)
		return counters, err
	}
	log.Println("CountersForLastMonth: ", counters)
	return counters, nil
}
