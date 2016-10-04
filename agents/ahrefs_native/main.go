package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/kardianos/osext"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/maddevsio/screen-monitoring/agents/ahrefs_native/service"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

type Settings struct {
	ID      string `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Content string `json:"content"`
}

func AgentRegistration(url, project_name, data_metrics string) {
	s := Settings{ID: "native_ahrefs_agent", Width: 200, Height: 200}
	s.Content = "<div><p>Native Ahrefs data mertrics of " + project_name + "</p><h1>" + data_metrics + "</h1></div>"
	log.Printf("Register to: %s", url)
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
	log.Println("response data_metrics:", data_metrics)
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
	log.Printf("%+v", metrics_data)
	metrics_data.CurrentRanges = metrics_data.CurrentRanges[:len(metrics_data.CurrentRanges)-1]
	metrics_data.MovementRanges = append(metrics_data.MovementRanges[:0], metrics_data.MovementRanges[0+1:]...)
	log.Printf("%+v", metrics_data)
	t := template.New("Metrics Data")
	folderPath, err := osext.ExecutableFolder()
	t, err = t.Parse("<table><tr><th>Organic keywords</th><th>Tracked keywords</th></tr>" +
		"<tr><td><b>All</b>: {{.All.Formated}} {{.All.Delta}}</td><td rowspan=2 colspan=2 style='text-align:center;'><span>{{.Keywords_tracked}} &uarr;{{.MovementTotal.Up}} &darr;{{.MovementTotal.Down}}</span></td></tr>" +
		"<tr><td><b>Us</b>: {{.Us.Formated}} {{.Us.Delta}}</td><td></td></tr>" +
		"<tr><td><b>Uk</b>: {{.Uk.Formated}} {{.Uk.Delta}}</td><td rowspan=4>" +
		"<ul style='list-style-type: none; display: inline-block; margin: 0; padding: 0 10px 0 0;'><li># 1-3</li><li># 4-10</li><li># 11-20</li><li># 21-50</li></ul>" +
		"<ul style='list-style-type: none; display: inline-block; margin: 0; padding: 0 10px 0 0;'>{{range .CurrentRanges}}<li>{{.}}</li>{{end}}</ul>" +
		"<ul style='list-style-type: none; display: inline-block; margin: 0; padding: 0; text-align=right;'>{{range .MovementRanges}}<li>{{.}}</li>{{end}}</ul>" +
		"</td></tr>" +
		"<tr><td><b>Au</b>: {{.Au.Formated}} {{.Au.Delta}}</td></tr>" +
		"<tr><td><b>Ca</b>: {{.Ca.Formated}} {{.Ca.Delta}}</td></tr>" +
		"<tr><td><b>Us</b>: {{.Us.Formated}} {{.Us.Delta}}</td></tr></table>")
	if err != nil {
		log.Println(err)
	}
	f, err := os.Create(folderPath + "/index.html")
	t.Execute(f, metrics_data)
	f.Close()

	AgentRegistration(*dashboardURL, *ahrefsProject, "test")
	e := echo.New()
	e.File("/", "index.html")
	e.Run(standard.New(*httpAddr))
}

type GetMetricsDataRequest struct {
	EMAIL    string `json:"ahrefsEmail"`
	PASSWORD string `json:"ahrefsPassword"`
	PROJECT  string `json:"ahrefsProject"`
}

type GetMetricsDataResponse struct {
	OrganicKeywords []byte `json:"organic_keywords"`
	TrackedKeywords []byte `json:"tracked_keywords"`
	Error           error  `json:"error"`
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
