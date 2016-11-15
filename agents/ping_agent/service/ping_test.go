package service

import (
	"testing"
	"strings"
)

func TestPingCommand(t *testing.T) {
	svc := agentService{}
	result := svc.CheckResponseTime("google.com")

	if !strings.HasSuffix(result, " ms")  {
		t.Errorf("Got empty result. Expected, got: %s", result)
	}
}
