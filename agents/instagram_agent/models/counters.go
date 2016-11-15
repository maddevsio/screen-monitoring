package models

import (
	"log"
	"time"

	"upper.io/db.v2/lib/sqlbuilder"
)

type Counter struct {
	Created    time.Time `db:"created" json:"created"`
	Username   string    `db:"username" json:"username"`
	Media      int64     `db:"media" json:"media"`
	Follows    int64     `db:"follows" json:"follows"`
	FollowedBy int64     `db:"followed_by" json:"followed_by"`
}

type AverageCounter struct {
	Date       string  `db:"date"`
	Media      float64 `db:"media"`
	Follows    float64 `db:"follows"`
	FollowedBy float64 `db:"followed_by"`
}

type CounterObject struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type CountersLastMonthResponse struct {
	Media      []CounterObject `json:"media"`
	Follows    []CounterObject `json:"follows"`
	FollowedBy []CounterObject `json:"followed_by"`
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

func (db *DB) CountersFindLast() (*Counter, error) {
	var counter *Counter
	q := db.SelectFrom("counters").OrderBy("created DESC").Limit(1)
	err := q.One(&counter)
	if err != nil {
		return nil, err
	}
	return counter, nil
}

func (db *DB) CountersLastMonth() ([]*AverageCounter, error) {
	var avgCounters []*AverageCounter

	rows, err := db.Query(`
SELECT
  	strftime('%Y-%m-%d', created) as date,
    avg(media) as media,
    avg(follows) as follows,
    avg(followed_by) as followed_by
FROM
    counters
WHERE
    created > datetime('now', '-1 month')
GROUP BY
    strftime('%Y-%m-%d', created)`)
	if err != nil {
		log.Fatal("Err getCountersLastMonth", err)
		return avgCounters, err
	}
	iter := sqlbuilder.NewIterator(rows)
	err = iter.All(&avgCounters)
	if err != nil {
		log.Fatal("Err cant iterrate", err)
		return nil, err
	}
	log.Println("CountersLastMonth: ", avgCounters)
	return avgCounters, nil
}
