package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}


func NewRouter(routes []*Route) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, r := range routes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(r.HandlerFunc)
	}

	return router
}