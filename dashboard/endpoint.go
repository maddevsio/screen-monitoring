package dashboard

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func makePagesEndpoint(svc DashboardService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		v, err := svc.GetPages()
		if err != nil {
			return []Widget{}, nil
		}
		return v, nil
	}
}

func makeRegisterWidgetEndpoint(svc DashboardService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Widget)
		v, err := svc.Register(req)
		if err != nil {
			return v, err
		}
		return v, nil
	}
}

func makeRegisterWidgetToPageEndpoint(svc DashboardService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(struct {
			pageId   int64
			widgetId string
		})
		v, err := svc.RegisterToPage(req.pageId, req.widgetId)
		if err != nil {
			return v, err
		}
		return v, nil
	}
}

func decodeRegisterWidgetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Widget
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeRegisterWidgetToPageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	pageId, _ := strconv.ParseInt(vars["pageId"], 10, 64)
	widgetId := vars["widgetId"]

	var request = struct {
		pageId   int64
		widgetId string
	}{pageId: pageId, widgetId: widgetId}

	return request, nil
}

func decodePagesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
