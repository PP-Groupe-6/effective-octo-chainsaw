package account_microservice

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(s AccountService, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeAccountEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET		/amount/ 		returns the amount of the param account
	// GET		/users/ 		returns the informations of the param account
	// POST 	/users/			adds a user

	r.Methods("GET").Path("/amount/").Handler(httptransport.NewServer(
		e.GetAmountEndpoint,
		decodeAmountRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/users/").Handler(httptransport.NewServer(
		e.GetUserInformationEndpoint,
		decodeUserInformationRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/users/").Handler(httptransport.NewServer(
		e.AddEndpoint,
		decodeAddRequest,
		encodeResponse,
		options...,
	))

	return r
}

type errorer interface {
	error() error
}

func decodeAmountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req GetAmountRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeAddRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req AddRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeUserInformationRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req GetUserInformationRequest

	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrNotAnId, ErrNotFound:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
