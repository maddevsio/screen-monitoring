package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/maddevsio/screen-monitoring/agents/ahrefs_native/service"
	"github.com/stretchr/testify/assert"
)

func TestAgentRegistration(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://127.0.0.1:8080/dashboard/v1/register",
		func(req *http.Request) (*http.Response, error) {
			agent_info := make(map[string]interface{})
			if err := json.NewDecoder(req.Body).Decode(&agent_info); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, err := httpmock.NewJsonResponse(200, "{'Success':true}")
			fmt.Println(resp)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
}

type mockedAhrefsServiceInterface interface {
	GetMetricsData() (metrics_data *service.MetricsData, err error)
}

type mockedAhrefsService struct{}

type MovementTotal struct {
	Up   int `json:"up"`
	Down int `json:"down"`
}

func (*mockedAhrefsService) GetMetricsData() (*service.MetricsData, error) {
	country := &service.Country{
		Formated: "1",
		Delta:    "2",
	}
	mData := &service.MetricsData{
		OrganicKeywords: service.OrganicKeywords{
			All: *country,
			Us:  *country,
			Uk:  *country,
			Au:  *country,
			Ca:  *country,
		},
		MovementRanges:   []int{2, 3, 5, 7, 11, 13},
		CurrentRanges:    []int{2, 3, 5, 7, 11, 13},
		Keywords_tracked: 1,
		MovementTotal: service.MovementTotal{
			Up:   1,
			Down: 2,
		},
	}
	return mData, nil
}

func TestShowData(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req := new(http.Request)
	c := e.NewContext(standard.NewRequest(req, e.Logger()), standard.NewResponse(rec, e.Logger()))
	c.SetPath("/data")
	metricsJSON := `{"organic_keywords":{"all":{"formated":"1","delta":"2"},"us":{"formated":"1","delta":"2"},"uk":{"formated":"1","delta":"2"},"au":{"formated":"1","delta":"2"},"ca":{"formated":"1","delta":"2"}},"movementRanges":[2,3,5,7,11,13],"currentRanges":[2,3,5,7,11,13],"keywords_tracked":1,"movementTotal":{"up":1,"down":2}}`
	env := Env{svc: &mockedAhrefsService{}}
	if assert.NoError(t, env.showData(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, metricsJSON, rec.Body.String())
	}
}
