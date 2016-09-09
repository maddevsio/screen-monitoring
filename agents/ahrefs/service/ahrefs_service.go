package service

import (
	"fmt"
	"github.com/ddliu/go-httpclient"
	"golang.org/x/net/html"
	"log"
)

type AhrefsService interface {
	SignIn(string, string) string
}

type ahrefsService struct{}

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

func (ahrefsService) SignIn(email, password string) string {
	response, _ := httpclient.Defaults(httpclient.Map{httpclient.OPT_DEBUG: true}).
		WithHeader("Accept-Language", "en-us").
		WithHeader("Referer", "https://ahrefs.com/").
		Get("https://ahrefs.com/", nil)

	root, _ := html.Parse(response.Body)
	element, _ := getElementByName("_token", root)
	token := ""
	for _, a := range element.Attr {
		if a.Key == "value" {
			token = a.Val
		}
	}
	fmt.Println("GET Request: ")
	log.Println(response.Request)
	fmt.Println("GET Response: ")
	log.Println(response.Response)

	response, err := httpclient.Defaults(httpclient.Map{httpclient.OPT_DEBUG: true}).
		WithHeader(":authority", "ahrefs.com").
		//Causes stream error: stream ID 3; PROTOCOL_ERROR
		//WithHeader(":method", "POST").
		//WithHeader(":path", "/user/login").
		//WithHeader(":scheme", "https").
		WithHeader("Accept", "*/*").
		WithHeader("Accept-Encoding", "gzip, deflate, br").
		WithHeader("Accept-Language", "en-US,en;q=0.8,ru;q=0.6").
		WithHeader("cache-control", "no-cache").
		WithHeader("content-type", "application/x-www-form-urlencoded").
		WithHeader("origin", "https://ahrefs.com").
		WithHeader("pragma", "no-cache").
		WithHeader("x-csrf-token", token).
		WithHeader("referer", "https://ahrefs.com/user/login").
		WithHeader("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36").
		WithHeader("x-requested-with", "XMLHttpRequest").
		//WithCookie(&http.Cookie{
		//	Name: "XSRF-TOKEN",
		//	Value: csrf_token,
		//}).
		//	WithCookie(&http.Cookie{
		//	Name: "ahrefs_cookie",
		//	Value: ahrefs_cookie,
		//}).
		Post("https://ahrefs.com/user/login/", map[string]string{
			"email":     email,
			"password":  password,
			"_token":    token,
			"return_to": "/",
		})
	fmt.Println("*** Cookies POST: ")
	for k, v := range httpclient.CookieValues("https://ahrefs.com/") {
		fmt.Printf("*%s : %s\n", k, v)
	}
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Request: ")
	log.Println(response.Request)
	fmt.Println("Response: ")
	log.Println(response.Response)

	response, err = httpclient.Defaults(httpclient.Map{httpclient.OPT_DEBUG: true}).
		WithHeader(":authority", "ahrefs.com").
		//Causes stream error: stream ID 5; PROTOCOL_ERROR
		//WithHeader(":method", "GET").
		//WithHeader(":path", "/").
		//WithHeader(":scheme", "https").
		WithHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8").
		WithHeader("Accept-Encoding", "gzip, deflate, sdch, br").
		WithHeader("Accept-Language", "en-US,en;q=0.8,ru;q=0.6").
		WithHeader("cache-control", "no-cache").
		WithHeader("referer", "https://ahrefs.com/user/login").
		WithHeader("pragma", "no-cache").
		WithHeader("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36").
		//WithCookie(&http.Cookie{
		//	Name: "XSRF-TOKEN",
		//	Value: csrf_token,
		//}).
		//	WithCookie(&http.Cookie{
		//	Name: "ahrefs_cookie",
		//	Value: ahrefs_cookie,
		//}).
		Get("https://ahrefs.com/dashboard/metrics", nil)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Request: ")
	log.Println(response.Request)
	fmt.Println("Response: ")
	log.Println(response.Response)

	defer response.Body.Close()

	return string(response.StatusCode)
}

func NewService() ahrefsService {
	return ahrefsService{}
}
