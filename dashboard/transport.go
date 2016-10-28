package dashboard

import (
	"encoding/json"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func MakeHandler(ctx context.Context, svc DashboardService, logger kitlog.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}
	renderDashboardHandler := kithttp.NewServer(
		ctx,
		makePagesEndpoint(svc),
		decodeNoParamsRequest,
		encodeResponse,
		options...,
	)

	unregisteredWidgetsHandler := kithttp.NewServer(
		ctx,
		makeUnregisteredWidgetsEndpoint(svc),
		decodeNoParamsRequest,
		encodeResponse,
		options...,
	)

	widgetRegisterHandler := kithttp.NewServer(
		ctx,
		makeRegisterWidgetEndpoint(svc),
		decodeRegisterWidgetRequest,
		encodeResponse,
		options...,
	)

	widgetRegisterToPageHandler := kithttp.NewServer(
		ctx,
		makeRegisterWidgetToPageEndpoint(svc),
		decodeRegisterWidgetToPageRequest,
		encodeResponse,
		options...,
	)

	pageInsertHandler := kithttp.NewServer(
		ctx,
		makeInsertPageEndpoint(svc),
		decodeInsertPageRequest,
		encodeResponse,
		options...,
	)

	r := mux.NewRouter()

	r.Handle("/dashboard/v1/pages",
		renderDashboardHandler).Methods("GET")
	r.Handle("/dashboard/v1/register",
		widgetRegisterHandler).Methods("POST")
	r.Handle(`/dashboard/v1/register/{widgetId:\w+}/page/{pageId:\d+}`,
		widgetRegisterToPageHandler).Methods("GET")
	r.Handle("/dashboard/v1/pages/new",
		pageInsertHandler).Methods("POST")
	r.Handle("/dashboard/v1/widgets/unregistered",
		unregisteredWidgetsHandler).Methods("GET")
	return r
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
