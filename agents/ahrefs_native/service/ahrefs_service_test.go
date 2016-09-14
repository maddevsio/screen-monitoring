package service

import (
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	var (
		email    = envString("AHREFS_EMAIL", "test@mail.com")
		password = envString("AHREFS_PASSWORD", "password")
	)

	svc := ahrefsService{}
	result := svc.SignIn(email, password)

	if result != nil {
		t.Errorf("Expected to be true. Got %s", result)
	}
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
