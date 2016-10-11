package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/maddevsio/screen-monitoring/agents/instagram_agent/models"
	"github.com/stretchr/testify/assert"
)

type mockDB struct{}

func (mdb *mockDB) CountersCreate(*models.Counter) error {
	return nil
}

func (mdb *mockDB) CountersFindLast() (*models.Counter, error) {
	counter := &models.Counter{time.Date(2016, time.October, 21, 0, 0, 0, 0, time.Local), "testuser", 10, 15, 20}
	return counter, nil
}

func (mdb *mockDB) CountersLastMonth() ([]*models.AverageCounter, error) {
	avgCounters := make([]*models.AverageCounter, 0)
	avgCounters = append(avgCounters, &models.AverageCounter{"2016-10-04", 10, 15, 20})
	avgCounters = append(avgCounters, &models.AverageCounter{"2016-10-05", 12, 15, 20})
	return avgCounters, nil
}

func TestCountersLast(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req := new(http.Request)
	c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
	c.SetPath("/counters")
	env := Env{db: &mockDB{}}
	countersJSON := `{"created":"2016-10-21T00:00:00+06:00","username":"testuser","media":10,"follows":15,"followed_by":20}`

	// Assertions
	if assert.NoError(t, env.countersLast(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, countersJSON, rec.Body.String())
	}
}
