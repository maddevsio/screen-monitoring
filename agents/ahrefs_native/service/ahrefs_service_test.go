package service

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var (
	project = "http://myproject.com"
)

func TestGetHash(t *testing.T) {
	body, err := ioutil.ReadFile("mocked_dashboard_response.html")
	if err != nil {
		fmt.Println(err)
	}
	_, status := getHash(body, project)
	if status == false {
		t.Errorf("Error, hash not found!")
	}
}

func TestGetToken(t *testing.T) {
	body, err := ioutil.ReadFile("mocked_home_response.html")
	if err != nil {
		fmt.Println(err)
	}
	token := getToken(body)
	if token == "" {
		t.Errorf("Error, token not found!")
	}
}
