package car

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	router := mux.NewRouter()
	endpoints := MakeServerEndpoints(s)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST    /cars/
	// GET     /cars/:id
	// PUT     /cars/:id
	// DELETE  /cars/:id

	router.Methods("POST").Path("/cars/").Handler(httptransport.NewServer(
		endpoints.PostCarEndpoint,
		decodePostCarRequest,
		encodeResponse,
		options...,
	))

	router.Methods("GET").Path("/cars/{id}").Handler(httptransport.NewServer(
		endpoints.GetCarEndpoint,
		decodeGetCarRequest,
		encodeResponse,
		options...,
	))

	router.Methods("PUT").Path("/cars/{id}").Handler(httptransport.NewServer(
		endpoints.PutCarEndpoint,
		decodePutCarRequest,
		encodeResponse,
		options...
	))

	router.Methods("DELETE").Path("/cars/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteCarEndpoint,
		decodeDeleteCarRequest,
		encodeResponse,
		options...,
	))

	return router
}

func decodePostCarRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req PostCarRequest

	err := json.NewDecoder(r.Body).Decode(&req.Car)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetCarRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	idRaw, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	id, e := strconv.Atoi(idRaw)

	if e != nil {
		return nil, ErrBadRouting
	}

	return GetCarRequest{Id: id}, nil
}

func decodePutCarRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	idRaw, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	id, err := strconv.Atoi(idRaw)

	if err != nil {
		return nil, ErrBadRouting
	}

	var car Car

	err = json.NewDecoder(r.Body).Decode(&car)

	if err != nil {
		return nil, err
	}

	return PutCarRequest{
		Id:  id,
		Car: car,
	}, nil
}

func decodeDeleteCarRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	idRaw, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}

	id, err := strconv.Atoi(idRaw)

	if err != nil {
		return nil, ErrBadRouting
	}

	return DelCarRequest{Id: id}, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(errorer)
	if ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
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
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
