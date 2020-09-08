package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"microservice_gokit_base/src/application/endpoints"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPOrder wires Go kit endpoints to the HTTP transport.
func NewHTTPOrder(
	svcEndpoints endpoints.IOrderEndpoints,
	logger log.Logger, baseURL string,
) http.Handler {
	// set-up router and initialize http endpoints
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}
	// HTTP Post - /orders
	r.Methods("POST").Path(baseURL + "orders").Handler(kithttp.NewServer(
		svcEndpoints.CreateEndpoint(),
		decodeCreateRequest(),
		encodeResponse,
		options...,
	))
	// HTTP Get - /orders/status
	r.Methods("GET").Path(baseURL + "orders/count").Handler(kithttp.NewServer(
		svcEndpoints.CountEndpoint(),
		decodeCount,
		encodeResponse,
		options...,
	))
	// HTTP Post - /orders/{id}
	r.Methods("GET").Path(baseURL + "orders/id/{id}").Handler(kithttp.NewServer(
		svcEndpoints.GetByIDEndpoint(),
		decodeGetByIDRequest,
		encodeResponse,
		options...,
	))

	// HTTP Get - /orders
	r.Methods("GET").Path(baseURL + "orders").Handler(kithttp.NewServer(
		svcEndpoints.GetAllEndpoint(),
		decodeGetAll,
		encodeResponse,
		options...,
	))

	// HTTP Get - /orders
	r.Methods("GET").Path(baseURL + "orders/saludo").Handler(kithttp.NewServer(
		svcEndpoints.SaludoEndpoint(),
		decodeGetAll,
		encodeResponse,
		options...,
	))

	// HTTP Put - /orders/status
	r.Methods("PUT").Path(baseURL + "orders/status").Handler(kithttp.NewServer(
		svcEndpoints.ChangeStatusEndpoint(),
		decodeChangeStausRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeCreateRequest() kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (request interface{}, err error) {
		var req endpoints.CreateRequest
		if e := json.NewDecoder(r.Body).Decode(&req.Order); e != nil {
			return nil, ErrBadRequest(e)
		}
		return req, nil
	}
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting(errors.New("No id parameter found"))
	}
	return endpoints.GetByIDRequest{ID: id}, nil
}

func decodeChangeStausRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req endpoints.ChangeStatusRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, ErrBadRequest(e)
	}
	return req, nil
}

func decodeGetAll(_ context.Context, r *http.Request) (request interface{}, err error) {
	page, _ := strconv.Atoi(r.FormValue("page"))
	size, _ := strconv.Atoi(r.FormValue("size"))

	if size != 0 {
		return endpoints.GetAllRequest{Page: int64(page), Size: int64(size)}, nil
	}
	return endpoints.GetAllRequest{Page: 0, Size: 0}, nil
}

func decodeCount(_ context.Context, r *http.Request) (request interface{}, err error) {
	return endpoints.CountRequest{}, nil
}
