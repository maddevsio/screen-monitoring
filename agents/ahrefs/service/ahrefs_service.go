package service

import (
	"golang.org/x/net/html"
	curl "github.com/andelf/go-curl"
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"fmt"
	"strings"
)

type AhrefsService interface {
	SignIn(string, string) string
}

type ahrefsService struct{}


func (ahrefsService) SignInAndGetDashboard(email, password string, verbose bool) string {
	easy := curl.EasyInit()
	token := ""
	receivedHTML := ""
	defer easy.Cleanup()

	getContent := func(body []byte, userdata interface{}) bool {
		receivedHTML += string(body)
		data, exists := getToken(body)
		if exists {
			token = data
		}

		return true
	}

	//first call
	easy.Setopt(curl.OPT_URL, "https://ahrefs.com/user/login")
	easy.Setopt(curl.OPT_SSL_VERIFYPEER, 1)
	easy.Setopt(curl.OPT_WRITEFUNCTION, getContent)
	easy.Setopt(curl.OPT_USERAGENT, "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko")
	easy.Setopt(curl.OPT_VERBOSE, verbose)
	easy.Setopt(curl.OPT_FOLLOWLOCATION, 1)
	easy.Setopt(curl.OPT_COOKIEJAR, "./cookiejar")
	easy.Setopt(curl.OPT_COOKIEFILE, "./cookiejar")
	easy.Setopt(curl.OPT_NOBODY, 0)
	easy.Perform()

	// lame and stupid method, but we have only callback for receiving data from curl_easy
	if strings.Contains(receivedHTML, "<strong>Dashboard") {
		return receivedHTML
	}

	receivedHTML = ""

	//second call in we need it (after first call we can be in Dashboard, thanks to cookieJar)
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

	return receivedHTML
}

func getToken(body []byte) (string, bool) {
	rootNode, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return "", false
	}

	meta := goquery.NewDocumentFromNode(rootNode).Find("meta[name=_token]")
	token, exists := meta.Attr("content")
	if exists == false {
		return "", false
	}

	return token, true

}

func NewService() ahrefsService {
	return ahrefsService{}
}
