package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/maddevsio/screen-monitoring/agents/ping_agent/service"
)

type Settings struct {
	ID      string `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Content string `json:"content"`
}

func AgentRegistration(url, hostname, gauge string) {
	s := Settings{ID: "ping_agent", Width: 200, Height: 200}
	s.Content = "<div><p>ping " + hostname + "</p><h1>" + gauge + "</h1></div>"
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
	log.Println("response Time:", gauge)
	body, _ := ioutil.ReadAll(client_resp.Body)
	log.Println("response Body:", string(body))
}

const (
	defaultPort         = "8090"
	defaultDashboardURL = "http://localhost:8080/dashboard/v1/register"
	defaultTargetHost   = "google.com"
)

func main() {
	var (
		addr         = envString("PORT", defaultPort)
		dashboardUrl = envString("DASHBOARD_URL", defaultDashboardURL)
		targetHost   = envString("TARGET_HOST", defaultTargetHost)

		httpAddr     = flag.String("httpAddr", ":"+addr, "HTTP listen address")
		dashboardURL = flag.String("dashboardURL", dashboardUrl, "Dashboard service URL")
		hostName     = flag.String("targetHost", targetHost, "Target hostname and port")

		ctx = context.Background()
		svc = service.NewService()
	)
	flag.Parse()

	timeTotal := svc.CheckResponseTime(*hostName)
	AgentRegistration(*dashboardURL, *hostName, timeTotal)

	checkResponseTimeHandler := httptransport.NewServer(
		ctx,
		makeCheckResponseTimeEndpoint(svc),
		decodeCheckResponseTimeRequest,
		encodeResponse,
	)
	http.Handle("/check", checkResponseTimeHandler)
	log.Printf("Http listen address: %s", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

func makeCheckResponseTimeEndpoint(svc service.AgentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(checkResponseTimeRequest)
		timeTotal := svc.CheckResponseTime(req.URL)
		return checkResponseTimeResponse{timeTotal}, nil
	}
}

func decodeCheckResponseTimeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request checkResponseTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type checkResponseTimeRequest struct {
	URL string `json:"url"`
}

type checkResponseTimeResponse struct {
	ResponseTimeTotal string `json:"response_time_total"`
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
