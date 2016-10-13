package service

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestGetMetricsData(t *testing.T) {
	var (
		email    = envString("AHREFS_EMAIL", "test@mail.com")
		password = envString("AHREFS_PASSWORD", "password")
		project  = envString("AHREFS_PROJECT", "http://projectname.com")
	)

	svc := ahrefsService{}
	organic_keywords, _, err := svc.GetMetricsData(email, password, project)
	fmt.Println("hello")

	fmt.Println(reflect.TypeOf(organic_keywords))
	if reflect.TypeOf(organic_keywords) != nil {
		t.Errorf("Error, expected byte arrays. Got %s", err)
	}
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
