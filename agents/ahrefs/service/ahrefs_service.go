package service

import (
	"fmt"
	"github.com/ddliu/go-httpclient"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/http/httputil"
	curl "github.com/andelf/go-curl"
	"bytes"
	"net/url"
)

type AhrefsService interface {
	SignIn(string, string) string
}

type ahrefsService struct{}


func (ahrefsService) SignIn(email, password string) string {
	easy := curl.EasyInit()
	token := ""
	defer easy.Cleanup()

	fooTest := func(body []byte, userdata interface{}) bool {
		//fmt.Print(string(body))
		root, _ := html.Parse(bytes.NewReader(body))
		element, _ := getElementByName("_token", root)
		if element != nil {
			for _, a := range element.Attr {
				if a.Key == "value" {
					token = a.Val
				}
			}
			//fmt.Print("Token = " + token + "\n")
		}
		return true
	}

	//first call
	easy.Setopt(curl.OPT_URL, "https://ahrefs.com/user/login")
	easy.Setopt(curl.OPT_SSL_VERIFYPEER, 1)
	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)
	easy.Setopt(curl.OPT_USERAGENT, "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko")
	easy.Setopt(curl.OPT_VERBOSE, 1)
	easy.Setopt(curl.OPT_FOLLOWLOCATION, 1)
	easy.Setopt(curl.OPT_COOKIEJAR, "/tmp/3")
	easy.Setopt(curl.OPT_COOKIEFILE, "/tmp/3")
	easy.Setopt(curl.OPT_NOBODY, 0)
	easy.Perform()

	//second call
	easy.Setopt(curl.OPT_URL, "https://ahrefs.com/user/login/")
	easy.Setopt(curl.OPT_HTTPHEADER, []string{
		"Referer: https://ahrefs.com/user/login",
		"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language: en-US,en;q=0.5",
		"Accept-Encoding: gzip, deflate",
		"Connection: keep-alive",
	})
	easy.Setopt(curl.OPT_POST, 1)
	form := url.Values{}
	form.Add("_token",    token)
	form.Add("email",     email)
	form.Add("password",  password)
	form.Add("return_to", "https://ahrefs.com/")
	postFields := form.Encode()
	fmt.Print(postFields + "\n")
	easy.Setopt(curl.OPT_POSTFIELDSIZE, len(postFields))
	easy.Setopt(curl.OPT_POSTFIELDS, postFields)
	easy.Perform()


	return "true"


	/*
		response, _ := httpclient.
			WithHeader("Accept-Language", "en-us").
			WithHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36").
			WithOption(httpclient.OPT_DEBUG, true).
			WithOption(httpclient.OPT_COOKIEJAR, true).
			WithOption(httpclient.OPT_FOLLOWLOCATION, false).
			Get("http://namba.kg/", nil)
	*/

	response, _ := httpclient.
		Defaults(httpclient.Map {
			httpclient.OPT_DEBUG: true,
			httpclient.OPT_USERAGENT: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36",
			"Accept-Language": "en-us",
			httpclient.OPT_COOKIEJAR: true,
			httpclient.OPT_FOLLOWLOCATION: false,
		}).Get("http://namba.kg/", nil)

	dump, _ := httputil.DumpResponse(response.Response, false)
	fmt.Printf("%s", dump)

	response2, _ := httpclient.Get("http://namba.kg/", nil)

	dump2, _ := httputil.DumpResponse(response2.Response, false)
	fmt.Printf("%s", dump2)

	return "true"

	csrf_token := ""
	ahrefs_cookie := ""
	fmt.Println("*** Cookies GET")
	for k, v := range httpclient.CookieValues("https://ahrefs.com/") {
		fmt.Printf("*%s : %s\n", k, v)
		if k == "XSRF-TOKEN" {
			csrf_token = v
		}
		if k == "ahrefs_cookie" {
			ahrefs_cookie = v
		}
	}

	root, _ := html.Parse(response.Body)
	element, _ := getElementByName("_token", root)
	token1 := ""
	for _, a := range element.Attr {
		if a.Key == "value" {
			token1 = a.Val
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
			"_token":    token1,
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
		WithCookie(&http.Cookie{
			Name:  "XSRF-TOKEN",
			Value: csrf_token,
		}).
		WithCookie(&http.Cookie{
			Name:  "ahrefs_cookie",
			Value: ahrefs_cookie,
		}).
		Get("https://ahrefs.com/dashboard/metrics", nil)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("Request: ")
	log.Println(response.Request)
	fmt.Println("Response: ")
	log.Println(response.Response)

	defer response.Body.Close()
	root, _ = html.Parse(response.Body)
	username, _ := getElementId("userAccountOptions", root)
	is_loggedin := "false"
	if username != nil {
		is_loggedin = "true"
	}
	log.Println(username)
	return string(is_loggedin)
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

