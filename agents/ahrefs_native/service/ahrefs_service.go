package service

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net/http/httputil"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type AhrefsService interface {
	GetMetricsData(email string, password string, project string) (organic_keywords []byte, tracked_keywords []byte, err error)
}

type ahrefsService struct{}

var client *http.Client

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
	res, err := client.Do(req)
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
	debug(httputil.DumpRequestOut(req, true))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else {
		defer res.Body.Close()
		debug(httputil.DumpResponse(res, true))
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

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

func (ahrefsService) GetMetricsData(email, password, project string) (organic_keywords []byte, tracked_keywords []byte, err error) {
	cJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, nil, err
	}

	client = &http.Client{
		Jar: cJar,
	}
	token := ""
	data, err := makeRequest("GET", "https://ahrefs.com/user/login", nil)
	if err != nil {
		return nil, nil, err
	}
	re := regexp.MustCompile(`.*<meta name="_token" content="(.*)" />.*`)
	token = re.FindStringSubmatch(string(data))[1]

	login_form := url.Values{}
	login_form.Add("_token", token)
	login_form.Add("email", email)
	login_form.Add("password", password)
	login_form.Add("return_to", "https://ahrefs.com/")
	dashboard_page, err := makeRequest("POST", "https://ahrefs.com/user/login", strings.NewReader(login_form.Encode()))
	if err != nil {
		return nil, nil, err
	}
	project_hash, _ := getHash(dashboard_page, project)
	organic_keywords_form := url.Values{}
	organic_keywords_form.Add("urls_hashes", project_hash)
	organic_keywords_form.Add("interval", "30")
	organic_keywords, err = makeRequestToGetMetricsData("POST", "https://ahrefs.com/dashboard/ajax/get/data",
		strings.NewReader(organic_keywords_form.Encode()), token)
	if err != nil {
		return nil, nil, err
	}
	tracked_keywords_form := url.Values{}
	tracked_keywords_form.Add("urls_hashes", project_hash)
	tracked_keywords, err = makeRequestToGetMetricsData("POST", "https://ahrefs.com/dashboard/ajax/get/rank-tracker",
		strings.NewReader(tracked_keywords_form.Encode()), token)
	if err != nil {
		return organic_keywords, nil, err
	}
	log.Println(string(organic_keywords))
	return organic_keywords, tracked_keywords, nil
}

func NewService() ahrefsService {
	return ahrefsService{}
}

func getElementByName(name string, n *html.Node) (element *html.Node, ok bool) {
	for _, a := range n.Attr {
		if a.Key == "name" && a.Val == name {
			return n, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = getElementByName(name, c); ok {
			return
		}
	}
	return
}

func getHash(body []byte, project_name string) (string, bool) {
	project_hash := ""
	rootNode, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return "Parse errror", false
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
	}
	return project_hash, true
}
