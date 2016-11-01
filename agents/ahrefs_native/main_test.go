package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestAgentRegistration(t *testing.T) {
	fmt.Println("start test")
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://127.0.0.1:8080/dashboard/v1/register",
		func(req *http.Request) (*http.Response, error) {
			agent_info := make(map[string]interface{})
			if err := json.NewDecoder(req.Body).Decode(&agent_info); err != nil {
				fmt.Println("err1")
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, err := httpmock.NewJsonResponse(200, "{'Success':true}")
			fmt.Println(resp)
			if err != nil {
				fmt.Println("err2")
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	fmt.Println("end test")
}

type mockedAhrefsService struct{}

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

func (*mockedAhrefsService) GetMetricsData() (*MetricsData, error) {
	country := Country{
		Formated: "1",
		Delta:    "2",
	}
	mData := &MetricsData{
		OrganicKeywords: OrganicKeywords{
			All: country,
			Us:  country,
			Uk:  country,
			Au:  country,
			Ca:  country,
		},
		MovementRanges:   []int{2, 3, 5, 7, 11, 13},
		CurrentRanges:    []int{2, 3, 5, 7, 11, 13},
		Keywords_tracked: 1,
		MovementTotal: MovementTotal{
			Up:   1,
			Down: 2,
		},
	}
	return mData, nil
}
