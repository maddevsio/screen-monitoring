package service

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
)

type AgentService interface {
	CheckResponseTime(string) string
}

type agentService struct{}

func (agentService) CheckResponseTime(url string) string {
	cmd := exec.Command("ping", "-c 1", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	data := out.String()
	re, _ := regexp.Compile(`\stime=(.*)\n`)
	info := re.FindStringSubmatch(data)

	return info[1]
}

func NewService() agentService {
	return agentService{}
}
