package dashboard

import (
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
	"net/http"
)

func makePagesEndpoint(svc DashboardService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		v, err := svc.GetPages()
		if err != nil {
			return nil, nil
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

func decodeRegisterWidgetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Widget
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodePagesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
