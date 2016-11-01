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

func TestLogin(t *testing.T) {
	dashboard, err := login("mytoken", "email", "password")
	fmt.Println(dashboard, err)
	if err != nil {
		t.Errorf("Error, can't login!")
	}
}

// type MyTestSuite struct {
// 	suite.Suite
// 	Client          *http.Client
// 	Server          *httptest.Server
// 	LastRequest     *http.Request
// 	LastRequestBody string
// 	ResponseFunc    func(http.ResponseWriter, *http.Request)
// }
//
// func (s *MyTestSuite) SetupSuite() {
// 	s.Server = httptest.NewTLSServer(
// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			body, _ := ioutil.ReadAll(r.Body)
// 			s.LastRequestBody = string(body)
// 			s.LastRequest = r
// 			if s.ResponseFunc != nil {
// 				s.ResponseFunc(w, r)
// 			}
// 		}))
// }
//
// func (s *MyTestSuite) TearDownSuite() {
// 	s.Server.Close()
// }
//
// func (s *MyTestSuite) SetupTest() {
// 	s.ResponseFunc = nil
// 	s.LastRequest = nil
// }
//
// func TestMySuite(t *testing.T) {
// 	suite.Run(t, new(MyTestSuite))
// }
//
// func (s *MyTestSuite) TestAgentRegistration() {
// 	//setup
// 	s.ResponseFunc = func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte(`{"id": 0, "name": "Hugh Jass"}`))
// 	}
//
// 	//test
// 	name := s.Client.AgentRegistration(0)
//
// 	//verify
// 	assert.Equal(s.T(), name, "Hugh Jass")
// 	assert.Equal(s.T(), "/people/0", s.LastRequest.URL.Path)
// }
