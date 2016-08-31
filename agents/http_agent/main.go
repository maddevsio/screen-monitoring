package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type AgentService interface {
	CheckResponseTime(string) (string, string)
}

type agentService struct{}

type Settings struct {
	ID      string `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Content string `json:"content"`
}

type HTTPAgent struct {
	AgentCheckURL string
	Hostname      string
	Gauge         string
}

func (agentService) CheckResponseTime(url string) (string, string) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		log.Println(err)
	}

	defer conn.Close()
	conn.Write([]byte("GET / HTTP/1.0\r\n\r\n"))

	startTime := time.Now()
	oneByte := make([]byte, 1)
	_, err = conn.Read(oneByte)
	if err != nil {
		log.Println(err)
	}
	timeStartTransfer := time.Since(startTime)

	_, err = ioutil.ReadAll(conn)
	if err != nil {
		log.Println(err)
	}
	elapsedTime := time.Since(startTime)

	return timeStartTransfer.String(), elapsedTime.String()
}

func AgentRegistration(url, hostname, gauge string) {
	s := Settings{ID: "http_agent", Width: 200, Height: 200}
	t, err := template.New("index.html").ParseFiles("tmpl/index.html")
	if err != nil {
		log.Fatalln(err)
	}
	agent := HTTPAgent{Hostname: hostname, Gauge: gauge}
	var b bytes.Buffer
	err = t.Execute(&b, agent)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(b.String())
	s.Content = b.String()
	log.Printf("Register dashboard service: %s", url)
	jsonStr, _ := json.Marshal(s)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("Connection failed: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
}

const (
	defaultPort         = "8090"
	defaultDashboardURL = "http://localhost:8080/dashboard/v1/register"
	defaultTarget       = "google.com:80"
)

func main() {
	var (
		addr  = envString("PORT", defaultPort)
		durl  = envString("DASHBOARD_URL", defaultDashboardURL)
		tHost = envString("TARGET_HOST", defaultTarget)

		httpAddr       = flag.String("httpAddr", ":"+addr, "HTTP listen address")
		dashboardURL   = flag.String("dashboardURL", durl, "Dashboard service URL")
		targetHostname = flag.String("targetHost", tHost, "Target hostname and port")

		ctx = context.Background()
		svc = agentService{}
	)

	flag.Parse()

	_, timeTotal := svc.CheckResponseTime(*targetHostname)
	AgentRegistration(*dashboardURL, *targetHostname, timeTotal)

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

func makeCheckResponseTimeEndpoint(svc AgentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(checkResponseTimeRequest)
		timeStart, timeTotal := svc.CheckResponseTime(req.URL)
		return checkResponseTimeResponse{timeStart, timeTotal}, nil
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
	TimeStartTransfer string `json:"response_time_start_transfer"`
	ResponseTimeTotal string `json:"response_time_total"`
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
