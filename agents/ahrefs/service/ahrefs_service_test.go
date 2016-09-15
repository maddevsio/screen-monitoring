package service

import (
	"flag"
	"os"
	"testing"
	"strings"
)

func init() {
	os.Remove("./cookiejar")
}

func TestLogin(t *testing.T) {
	var (
		email    = envString("AHREFS_EMAIL", "test@mail.com")
		password = envString("AHREFS_PASSWORD", "password")

		ahrefsEmail    = flag.String("ahrefsEmail", email, "Email address of your ahrefs.com account")
		ahrefsPassword = flag.String("ahrefsPassword", password, "Password")
	)

	svc := ahrefsService{}
	result := svc.SignInAndGetDashboard(*ahrefsEmail, *ahrefsPassword, true)

	if !strings.Contains(result, "<strong>Dashboard") {
		t.Error("Expected to be in Dashboard", nil)
	}
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
