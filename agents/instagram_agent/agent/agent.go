package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Settings struct {
	ID     string `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

func Register(dashboardURL string, indexURL string) error {
	s := Settings{
		ID:     "instagram_agent",
		Width:  600,
		Height: 900,
		URL:    indexURL,
	}

	jsonStr, _ := json.Marshal(s)
	req, err := http.NewRequest("POST", dashboardURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	clientResp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer clientResp.Body.Close()

	return nil
}
