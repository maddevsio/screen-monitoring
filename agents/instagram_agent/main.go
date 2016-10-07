package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/yanatan16/golang-instagram/instagram"

	"github.com/maddevsio/screen-monitoring/agents/instagram_agent/agent"
	"github.com/maddevsio/screen-monitoring/agents/instagram_agent/models"
)

const (
	defaultPort         = "8090"
	defaultDashboardURL = "http://localhost:8080/dashboard/v1/register"
)

var (
	clientID    = envString("CLIEND_ID", "")
	accessToken = envString("ACCESS_TOKEN", "")
)

type Env struct {
	db models.Datastore
}

func main() {
	var (
		addr         = envString("PORT", defaultPort)
		dashboardUrl = envString("DASHBOARD_URL", defaultDashboardURL)

		httpAddr     = flag.String("httpAddr", "127.0.0.1:"+addr, "HTTP listen address")
		dashboardURL = flag.String("dashboardURL", dashboardUrl, "Dashboard service URL")
	)

	flag.Parse()

	err := agent.Register(*dashboardURL, "http://"+*httpAddr)
	if err != nil {
		log.Fatal(err)
	}

	db, err := models.NewDB("instagram.db")
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{db}

	api := &instagram.Api{
		ClientId:    clientID,
		AccessToken: accessToken,
	}

	ticker := time.NewTicker(time.Minute * 2)
	go func() {
		for _ = range ticker.C {
			if ok, err := api.VerifyCredentials(); !ok {
				log.Fatal(err)
			}

			if resp, err := api.GetSelf(); err != nil {
				log.Fatal(err)
			} else {
				me := resp.User
				counter := &models.Counter{
					Username:   me.Username,
					Media:      me.Counts.Media,
					FollowedBy: me.Counts.FollowedBy,
					Follows:    me.Counts.Follows,
				}
				env.db.CountersCreate(counter)
			}
		}
	}()

	e := echo.New()
	e.File("/", "tmpl/index.html")
	e.Static("/static", "assets")
	e.GET("/counters", env.countersLast)
	e.GET("/counters-last-month", env.countersLastMonth)
	e.Run(standard.New(*httpAddr))
}

func (env *Env) countersLast(c echo.Context) error {
	counters, err := env.db.CountersFindLast()
	if err != nil {
		// TODO: handle this
		log.Fatal(err)
		return err
	}
	return c.JSON(http.StatusOK, counters)
}

func (env *Env) countersLastMonth(c echo.Context) error {
	avgCounters, err := env.db.CountersLastMonth()
	if err != nil {
		return err
	}

	var media []models.CounterObject
	var follows []models.CounterObject
	var followed_by []models.CounterObject

	for _, avgCounter := range avgCounters {
		media = append(media, models.CounterObject{
			Date:     avgCounter.Date,
			Counters: avgCounter.Media,
		})
		follows = append(follows, models.CounterObject{
			Date:     avgCounter.Date,
			Counters: avgCounter.Follows,
		})
		followed_by = append(followed_by, models.CounterObject{
			Date:     avgCounter.Date,
			Counters: avgCounter.FollowedBy,
		})
	}

	response := models.CountersLastMonthResponse{
		Media:      media,
		Follows:    follows,
		FollowedBy: followed_by,
	}

	return c.JSON(http.StatusOK, response)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
