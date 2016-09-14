package service

import (
	"fmt"
	"golang.org/x/net/html"
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
	easy.Setopt(curl.OPT_COOKIEJAR, "./cookiejar")
	easy.Setopt(curl.OPT_COOKIEFILE, "./cookiejar")
	easy.Setopt(curl.OPT_NOBODY, 0)
	easy.Perform()

	//second call
	easy.Setopt(curl.OPT_URL, "https://ahrefs.com/user/login")
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

