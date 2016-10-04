package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"upper.io/db.v2/lib/sqlbuilder"
	"upper.io/db.v2/sqlite"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/yanatan16/golang-instagram/instagram"
)

type Settings struct {
	ID     string `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type UserData struct {
	Username   string
	Media      int64
	Follows    int64
	FollowedBy int64
}

type Counters struct {
	Created    time.Time `db:"created" json:"created"`
	Media      int64     `db:"media" json:"media"`
	Follows    int64     `db:"follows" json:"follows"`
	FollowedBy int64     `db:"followed_by" json:"followed_by"`
}

type Service interface {
}

type service struct{}

func (service) Register(u UserData, url string, serverAddr string) {
	s := Settings{
		ID:     "instagram_agent",
		Width:  300,
		Height: 300,
		URL:    "http://" + serverAddr,
	}

	jsonStr, _ := json.Marshal(s)
	req, req_err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if req_err != nil {
		log.Println("Connection failed: ", req_err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client_resp, client_err := client.Do(req)
	if client_err != nil {
		log.Println("Connection failed: ", client_err)
	}
	defer client_resp.Body.Close()
}

func GetInstagramUserInfo() (u UserData, err error) {
	api := &instagram.Api{
		ClientId:    clientID,
		AccessToken: accessToken,
	}

	if ok, err := api.VerifyCredentials(); !ok {
		return UserData{}, err
	}

	if resp, err := api.GetSelf(); err != nil {
		return UserData{}, err
	} else {
		me := resp.User
		data := UserData{
			Username:   me.Username,
			Media:      me.Counts.Media,
			FollowedBy: me.Counts.FollowedBy,
			Follows:    me.Counts.Follows,
		}
		return data, nil
	}
}

func WriteData(d UserData, sess sqlbuilder.Database) {
	countersCollection := sess.Collection("counters")
	countersCollection.Insert(Counters{
		Created:    time.Now(),
		Media:      d.Media,
		Follows:    d.Follows,
		FollowedBy: d.FollowedBy,
	})

	res := countersCollection.Find()
	var counters []Counters
	err := res.All(&counters)
	if err != nil {
		log.Fatalf("res.All(): %q/n", err)
	}

	for _, counter := range counters {
		fmt.Printf("Created: %s, media: %d, follows: %d, followed_by: %d \n",
			counter.Created,
			counter.Media,
			counter.Follows,
			counter.FollowedBy,
		)
	}
	return
}

func GetCounters(sess sqlbuilder.Database) (Counters, error) {
	var counters Counters
	q := sess.SelectFrom("counters").OrderBy("created DESC").Limit(1)
	err := q.One(&counters)
	if err != nil {
		return Counters{}, err
	}
	return counters, nil
}

func GetCountersLastMonth(sess sqlbuilder.Database) {
	var counters []Counters
	q := sess.SelectFrom("counters").Where("created > datetime('created', '-1 month')")
	err := q.All(&counters)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(counters)
	return
}

const (
	defaultPort         = "8090"
	defaultDashboardURL = "http://localhost:8080/dashboard/v1/register"
)

var (
	clientID    = envString("CLIEND_ID", "")
	accessToken = envString("ACCESS_TOKEN", "")
)

func main() {
	var (
		addr         = envString("PORT", defaultPort)
		dashboardUrl = envString("DASHBOARD_URL", defaultDashboardURL)

		httpAddr     = flag.String("httpAddr", "127.0.0.1:"+addr, "HTTP listen address")
		dashboardURL = flag.String("dashboardURL", dashboardUrl, "Dashboard service URL")

		// ctx = context.Background()
		svc      = service{}
		settings = sqlite.ConnectionURL{
			Database: "instagram.db",
		}
	)
	flag.Parse()

	userData, err := GetInstagramUserInfo()
	if err != nil {
		log.Fatal(err)
	}

	svc.Register(userData, *dashboardURL, *httpAddr)

	sess, err := sqlite.Open(settings)
	if err != nil {
		log.Fatalf("db.Open(): %q\n", err)
	}
	defer sess.Close()

	GetCountersLastMonth(sess)

	ticker := time.NewTicker(time.Minute * 2)
	go func() {
		for _ = range ticker.C {
			userData, err := GetInstagramUserInfo()
			if err != nil {
				log.Fatal(err)
			}
			WriteData(userData, sess)
		}
	}()

	// time.Sleep(time.Minute * 2)
	// log.Println("Ticker stop")

	// pageSpeedScoreHandler := httptransport.NewServer(
	// 	ctx,
	// 	makePageSpeedScoreEndpoint(svc),
	// 	decodePageSpeedScoreRequest,
	// 	encodeResponse,
	// )
	// http.Handle("/check", pageSpeedScoreHandler)
	// log.Printf("Http listen address: %s", *httpAddr)
	// log.Fatal(http.ListenAndServe(*httpAddr, nil))

	// type User struct {
	// 	Name  string `json:"name" xml:"name" form:"name"`
	// 	Email string `json:"email" xml:"email" form:"email"`
	// }

	e := echo.New()
	e.File("/", "tmpl/index.html")
	e.GET("/counters", func(c echo.Context) error {
		// counters, err := GetCounters(sess)
		u := new(Counters)
		if err := c.Bind(u); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
	})
	e.Run(standard.New(*httpAddr))
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
