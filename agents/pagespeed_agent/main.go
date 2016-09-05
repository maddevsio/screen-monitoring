package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

type AgentService interface {
	PageSpeedScore(url string) *Score
}

type agentService struct{}

type Settings struct {
	ID      string `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Content string `json:"content"`
}

type Score struct {
	MobileSpeed     int
	MobileUsability int
	DesktopSpeed    int
}

type PageSpeedMobile struct {
	RuleGroups struct {
		Speed struct {
			Score int
		}
		Usability struct {
			Score int
		}
	}
}

type PageSpeedDesktop struct {
	RuleGroups struct {
		Speed struct {
			Score int
		}
	}
}

func (agentService) PageSpeedScore(url string) *Score {
	mURL := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v2/runPagespeed?url=%s&filter_third_party_resources=true&locale=ru&strategy=mobile&key=%s", url, apiKey)
	dURL := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v2/runPagespeed?url=%s&filter_third_party_resources=true&locale=ru&strategy=desktop&key=%s", url, apiKey)

	// Mobile
	res, err := http.Get(mURL)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(res.Body)
	var dataMobile PageSpeedMobile
	err = json.Unmarshal(body, &dataMobile)
	if err != nil {
		log.Fatal(err)
	}
	// Desktop
	res, err = http.Get(dURL)
	if err != nil {
		log.Fatal(err)
	}
	body, _ = ioutil.ReadAll(res.Body)
	var dataDesktop PageSpeedDesktop
	err = json.Unmarshal(body, &dataDesktop)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	s := &Score{
		MobileUsability: dataMobile.RuleGroups.Usability.Score,
		MobileSpeed:     dataMobile.RuleGroups.Speed.Score,
		DesktopSpeed:    dataDesktop.RuleGroups.Speed.Score,
	}
	return s
}

func AgentRegistration(url string, tURL string, sc Score) {
	s := Settings{ID: "pagespeed_agent", Width: 300, Height: 300}

	s.Content = "<div><p>PageSpeed: " + tURL + "</p>" +
		"<p " + getColor(sc.MobileSpeed) + ">MobileSpeed: " + strconv.Itoa(sc.MobileSpeed) + "</p>" +
		"<p " + getColor(sc.MobileUsability) + ">MobileUsability: " + strconv.Itoa(sc.MobileUsability) + "</p>" +
		"<p " + getColor(sc.DesktopSpeed) + ">DesktopSpeed: " + strconv.Itoa(sc.DesktopSpeed) + "</p></div>"

	log.Println(s.Content)

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
}

func getColor(v int) string {
	//return style string
	switch {
	case v >= 90:
		return `style="color:green;"`
	case v <= 89 && v >= 70:
		return `style="color:yellow;"`
	case v <= 69:
		return `style="color:red;"`
	}
	return ""
}

const (
	defaultPort         = "8090"
	defaultDashboardURL = "http://localhost:8080/dashboard/v1/register"
	defaultTargetUrl    = ""
)

var apiKey = envString("API_KEY", "")

func main() {

	var (
		addr         = envString("PORT", defaultPort)
		dashboardUrl = envString("DASHBOARD_URL", defaultDashboardURL)
		targetUrl    = envString("TARGET_URL", defaultTargetUrl)

		httpAddr     = flag.String("httpAddr", ":"+addr, "HTTP listen address")
		dashboardURL = flag.String("dashboardURL", dashboardUrl, "Dashboard service URL")
		tURL         = flag.String("targetURL", targetUrl, "Target URL")

		ctx = context.Background()
		svc = agentService{}
	)

	flag.Parse()

	score := svc.PageSpeedScore(*tURL)
	AgentRegistration(*dashboardURL, *tURL, *score)

	pageSpeedScoreHandler := httptransport.NewServer(
		ctx,
		makePageSpeedScoreEndpoint(svc),
		decodePageSpeedScoreRequest,
		encodeResponse,
	)
	http.Handle("/check", pageSpeedScoreHandler)
	log.Printf("Http listen address: %s", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

func makePageSpeedScoreEndpoint(svc AgentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pageSpeedScoreRequest)
		score := svc.PageSpeedScore(req.URL)
		return pageSpeedScoreResponse{MobileSpeed: score.MobileSpeed, MobileUsability: score.MobileUsability, DesktopSpeed: score.DesktopSpeed}, nil
	}
}

func decodePageSpeedScoreRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request pageSpeedScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type pageSpeedScoreRequest struct {
	URL string `json:"url"`
}

type pageSpeedScoreResponse struct {
	MobileSpeed     int `json:"mobile_speed"`
	MobileUsability int `json:"mobile_usability"`
	DesktopSpeed    int `json:"desktop_speed"`
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
