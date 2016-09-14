package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/maddevsio/screen-monitoring/agents/ahrefs/service"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Settings struct {
	ID      string `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Content string `json:"content"`
}

func AgentRegistration(url, project_name, data_metrics string) {
	s := Settings{ID: "ahrefs_agent", Width: 200, Height: 200}
	s.Content = "<div><p>Ahrefs data mertrics of " + project_name + "</p><h1>" + data_metrics + "</h1></div>"
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

		ctx = context.Background()
		svc = service.NewService()
	)

	flag.Parse()
	loginStatus := svc.SignInAndGetDashboard(*ahrefsEmail, *ahrefsPassword)
	log.Println(loginStatus)
	metrics_data := "No data" // Implement parsing data from ahrefs dashboard
	AgentRegistration(*dashboardURL, *ahrefsProject, metrics_data)
	LoginHandler := httptransport.NewServer(
		ctx,
		makeSignInEndpoint(svc),
		decodeSigInRequest,
		encodeResponse,
	)
	http.Handle("/check", LoginHandler)
	log.Printf("Http listen address: %s", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))

}

func makeSignInEndpoint(svc service.AhrefsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SigInRequest)
		loginStatus := svc.SignIn(req.EMAIL, req.PASSWORD)
		return SignInResponse{loginStatus}, nil
	}
}

func decodeSigInRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request SigInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type SigInRequest struct {
	EMAIL    string `json:"ahrefsEmail"`
	PASSWORD string `json:"ahrefsPassword"`
}

type SignInResponse struct {
	ResponseStatus string `json:"response_status"`
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
