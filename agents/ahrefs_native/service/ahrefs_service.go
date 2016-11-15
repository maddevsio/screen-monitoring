package service

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"bytes"
	"compress/gzip"
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

type AhrefsServiceStruct struct {
	email    string
	password string
	project  string
}

var cJar, err = cookiejar.New(nil)

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

var MyClient = &http.Client{
	Timeout: time.Second * 10,
	Jar:     cJar,
}

type AhrefsServiceInterface interface {
	GetMetricsData() (metrics_data *MetricsData, err error)
}

func makeRequest(method, url string, reader io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, reader)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "https://ahrefs.com/user/login")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko")
	res, err := MyClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 500 {
		return nil, errors.New("Internal Server Error")
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func makeRequestToGetMetricsData(method, url string, reader io.Reader, token string) (data []byte, err error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Referer", "https://ahrefs.com/dashboard/metrics")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8,ru;q=0.6")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("x-csrf-token", token)
	req.Header.Add("Content-Length", "56")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36")
	// debug(httputil.DumpRequestOut(req, true))
	res, err := MyClient.Do(req)
	if err != nil {
		return nil, err
	} else {
		defer res.Body.Close()
		// debug(httputil.DumpResponse(res, true))
		if res.StatusCode == 500 {
			return nil, errors.New("Internal Server Error")
		}
		switch res.Header.Get("Content-Encoding") {
		case "gzip":
			log.Println("GZIP")
			reader, err = gzip.NewReader(res.Body)
		default:
			reader = res.Body
		}
	}

	data, err = ioutil.ReadAll(reader)

	if err != nil {
		panic(err)
	}
	return data, nil
}

//Gets token from ahrefs main page
func getToken(data []byte) string {
	re := regexp.MustCompile(`.*<meta name="_token" content="(.*)" />.*`)
	token := re.FindStringSubmatch(string(data))[1]
	return token
}

//Logins to ahrefs and returns dashboard page in byte array
func login(token string, email string, password string) ([]byte, error) {
	login_form := url.Values{}
	login_form.Add("_token", token)
	login_form.Add("email", email)
	login_form.Add("password", password)
	login_form.Add("return_to", "https://ahrefs.com/")
	dashboard_page, err := makeRequest("POST", "https://ahrefs.com/user/login", strings.NewReader(login_form.Encode()))
	if err != nil {
		return nil, err
	}
	return dashboard_page, nil
}

//Sends xhr request to get organic keywords
func getOrganicKeyWords(token string, project_hash string) (organic_keywords []byte, err error) {
	organic_keywords_form := url.Values{}
	organic_keywords_form.Add("urls_hashes", project_hash)
	organic_keywords_form.Add("interval", "30")
	organic_keywords, err = makeRequestToGetMetricsData("POST", "https://ahrefs.com/dashboard/ajax/get/data",
		strings.NewReader(organic_keywords_form.Encode()), token)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return organic_keywords, nil
}

//Sends xhr request to get tracked keywords
func getTrackedKeyWords(token string, project_hash string) (tracked_keywords []byte, err error) {
	tracked_keywords_form := url.Values{}
	tracked_keywords_form.Add("urls_hashes", project_hash)
	tracked_keywords, err = makeRequestToGetMetricsData("POST", "https://ahrefs.com/dashboard/ajax/get/rank-tracker",
		strings.NewReader(tracked_keywords_form.Encode()), token)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return tracked_keywords, nil
}

func (a *AhrefsServiceStruct) GetMetricsData() (metrics_data *MetricsData, err error) {
	token := ""
	data, err := makeRequest("GET", "https://ahrefs.com/user/login", nil)
	if err != nil {
		return nil, err
	}

	token = getToken(data)

	dashboard_page, err := login(token, a.email, a.password)
	if err != nil {
		return nil, err
	}

	project_hash, _ := getHash(dashboard_page, a.project)

	organic_keywords, err := getOrganicKeyWords(token, project_hash)
	if err != nil {
		return nil, err
	}

	tracked_keywords, err := getTrackedKeyWords(token, project_hash)
	if err != nil {
		return nil, err
	}

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
	return metrics_data, nil
}

func NewService(email, password, project string) *AhrefsServiceStruct {
	return &AhrefsServiceStruct{email, password, project}
}

//Gets projects hash from ahrefs metrics. Takes dashboard page and projects name.
//Returns projects hash and status
func getHash(body []byte, project_name string) (string, bool) {
	project_hash := ""
	status := false
	rootNode, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return project_hash, status
	}
	doc := goquery.NewDocumentFromNode(rootNode)
	match_string := ""
	doc.Find("div.dashboard-media").Each(func(i int, s *goquery.Selection) {
		project_link_url, found := s.Find("a.c-green-primary").Attr("href")
		if found && project_link_url == project_name {
			log.Printf("Found %s project.", project_name)
			s.Find("th").Each(func(j int, th *goquery.Selection) {
				th_text := th.Text()
				if th_text == "Organic keywords " {
					match_string, _ = th.Html()
				}
			})
		}
	})
	re := regexp.MustCompile(`id="organic_keywords_alert_([a-f0-9]+)`)
	match_result := re.FindStringSubmatch(match_string)
	if match_result != nil {
		project_hash = match_result[1]
		status = true
	}
	return project_hash, status
}
