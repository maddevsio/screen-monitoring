package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/maddevsio/screen-monitoring/agents/ahrefs_native/service"
)

type Settings struct {
	ID      string `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Content string `json:"content"`
	Url     string `json:"url"`
}

func AgentRegistration(url, agent_url string) {
	s := Settings{ID: "native_ahrefs_agent", Width: 300, Height: 180, Content: "", Url: agent_url}
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

	log.Println("response Status:", client_resp.Status)
	body, _ := ioutil.ReadAll(client_resp.Body)
	log.Println("response Body:", string(body))
}

const (
	defaultPort         = "8090"
	defaultEmail        = "email@mail.com"
	defaultPassword     = "password"
	defaultProjectName  = "myproject.com"
	defaultDashboardURL = "http://localhost:8080/dashboard/v1/register"
)

type Country struct {
	Formated interface{} `json:"formated"`
	Delta    interface{} `json:"delta"`
}

type MetricsData struct {
	OrganicKeywords  `json:"organic_keywords"`
	MovementRanges   []int `json:"movementRanges"`
	CurrentRanges    []int `json:"currentRanges"`
	Keywords_tracked int   `json:"keywords_tracked"`
	MovementTotal    `json:"movementTotal"`
}

type OrganicKeywords struct {
	All Country `json:"all"`
	Us  Country `json:"us"`
	Uk  Country `json:"uk"`
	Au  Country `json:"au"`
	Ca  Country `json:"ca"`
}

type MovementTotal struct {
	Up   int `json:"up"`
	Down int `json:"down"`
}

func main() {

	var (
		addr         = envString("PORT", defaultPort)
		dashboardUrl = envString("DASHBOARD_URL", defaultDashboardURL)
		email        = envString("AHREFS_EMAIL", defaultEmail)
		password     = envString("AHREFS_PASSWORD", defaultPassword)
		project      = envString("AHREFS_PROJECT", defaultProjectName)

		httpAddr       = flag.String("httpAddr", ":"+addr, "HTTP listen address")
		dashboardURL   = flag.String("dashboardURL", dashboardUrl, "Dashboard service URL")
		ahrefsEmail    = flag.String("ahrefsEmail", email, "Email address of your ahrefs.com account")
		ahrefsPassword = flag.String("ahrefsPassword", password, "Password")
		ahrefsProject  = flag.String("ahrefsProject", project, "Name of the project which data metrics you"+
			" want to get. Be sure to use the exact name which is shown at ahrefs dahsboard.")

		svc = service.NewService()
	)

	flag.Parse()

	AgentRegistration(*dashboardURL, "http://localhost:8090/")
	e := echo.New()
	e.File("/", "index.html")
	e.GET("/data", func(c echo.Context) error {

		organic_keywords, tracked_keywords, err := svc.GetMetricsData(*ahrefsEmail, *ahrefsPassword, *ahrefsProject)
		if err != nil {
			log.Println(err)
		}
		metrics_data := MetricsData{}
		err = json.Unmarshal(organic_keywords, &metrics_data)
		if err != nil {
			log.Println("error: ", err)
		}

		err = json.Unmarshal(tracked_keywords, &metrics_data)
		if err != nil {
			log.Println("error: ", err)
		}
		metrics_data.CurrentRanges = metrics_data.CurrentRanges[:len(metrics_data.CurrentRanges)-1]
		metrics_data.MovementRanges = append(metrics_data.MovementRanges[:0], metrics_data.MovementRanges[0+1:]...)

		return c.JSON(http.StatusOK, metrics_data)
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
