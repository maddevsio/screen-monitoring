package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/maddevsio/screen-monitoring/dashboard"
	gometrics "github.com/rcrowley/go-metrics"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	logger := log.NewLogfmtLogger(os.Stderr)

	requestCount := gometrics.NewCounter()
	requestLatency := gometrics.NewHistogram(gometrics.NewUniformSample(100))
	countResult := gometrics.NewHistogram(gometrics.NewUniformSample(100))
	var svc dashboard.DashboardService
	svc = dashboard.NewDashboardService()
	svc = dashboard.NewLoggerService(logger, svc)
	svc = dashboard.NewInstrumentingMiddleware(
		requestCount,
		requestLatency,
		countResult,
		svc,
	)

	mux := http.NewServeMux()
	mux.Handle("/dashboard/v1/", dashboard.MakeHandler(ctx, svc, logger))
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)
	go func() {
		logger.Log("msg", "HTTP", "addr", ":8080")
		errs <- http.ListenAndServe(":8080", nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
