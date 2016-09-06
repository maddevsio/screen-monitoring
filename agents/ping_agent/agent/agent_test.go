package agent

import "testing"

func TestPingCommand(t *testing.T) {
	svc := agentService{}
	result := svc.CheckResponseTime("google.com")

	if result != "" {
		t.Errorf("Got empty result. Expected, got: %s", result)
	}
}
