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

func AgentRegistration(url, agent_url string) (bool, error) {
	s := Settings{ID: "native_ahrefs_agent", Width: 300, Height: 180, Content: "", Url: agent_url}
	jsonStr, _ := json.Marshal(s)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("AgentRegistration: ", err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client_resp, err := client.Do(req)
	if err != nil {
		log.Println("AgentRegistration: ", err)
		return false, err
	}
	defer client_resp.Body.Close()

	body, err := ioutil.ReadAll(client_resp.Body)
	if err != nil {
		log.Println("AgentRegistration: ", err)
		return false, err
	}
	log.Printf("Status: %s. Body: %s", client_resp.Status, string(body))
	return true, nil
}

const (
	defaultPort         = "8090"
	defaultEmail        = "email@mail.com"
	defaultPassword     = "password"
	defaultProjectName  = "myproject.com"
	defaultDashboardURL = "http://localhost:8080/dashboard/v1/register"
)

var (
	email    = envString("AHREFS_EMAIL", defaultEmail)
	password = envString("AHREFS_PASSWORD", defaultPassword)
	project  = envString("AHREFS_PROJECT", defaultProjectName)

	ahrefsEmail    = flag.String("ahrefsEmail", email, "Email address of your ahrefs.com account")
	ahrefsPassword = flag.String("ahrefsPassword", password, "Password")
	ahrefsProject  = flag.String("ahrefsProject", project, "Name of the project which data metrics you"+
		" want to get. Be sure to use the exact name which is shown at ahrefs dahsboard.")
)

type AhrefsService struct {
}

type Env struct {
	svc service.AhrefsServiceInterface
}

func main() {

	var (
		addr         = envString("PORT", defaultPort)
		dashboardUrl = envString("DASHBOARD_URL", defaultDashboardURL)
		httpAddr     = flag.String("httpAddr", ":"+addr, "HTTP listen address")
		dashboardURL = flag.String("dashboardURL", dashboardUrl, "Dashboard service URL")
	)

	flag.Parse()

	svc := service.NewService(*ahrefsEmail, *ahrefsPassword, *ahrefsProject)
	env := &Env{svc}

	_, err := AgentRegistration(*dashboardURL, "http://localhost:8090/")
	if err != nil {
		log.Println("Main: ", err)
	}

	e := echo.New()
	e.File("/", "index.html")
	e.GET("/data", env.showData)
	e.Run(standard.New(*httpAddr))
}

func (env *Env) showData(c echo.Context) error {
	metrics_data, err := env.svc.GetMetricsData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, metrics_data)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
