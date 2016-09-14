package service

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

type AhrefsService interface {
	SignIn(string, string) string
}

type ahrefsService struct{}

var client *http.Client

func makeRequest(method string, reader io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, "https://ahrefs.com/user/login", reader)

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
func (ahrefsService) SignIn(email, password string) error {
	cJar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return err
	}

	client = &http.Client{
		Jar: cJar,
	}
	token := ""
	data, err := makeRequest("GET", nil)
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`.*<meta name="_token" content="(.*)" />.*`)
	token = re.FindStringSubmatch(string(data))[1]
	fmt.Println(token)

	form := url.Values{}
	form.Add("_token", token)
	form.Add("email", email)
	form.Add("password", password)
	form.Add("return_to", "https://ahrefs.com/")
	data, err = makeRequest("POST", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	return nil
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

func getElementId(id string, n *html.Node) (element *html.Node, ok bool) {
	for _, a := range n.Attr {
		if a.Key == "id" && a.Val == id {
			return n, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = getElementByName(id, c); ok {
			return
		}
	}
	return
}
