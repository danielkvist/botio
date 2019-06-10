package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route represents a simple route for a mux.Router.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// NewRouter iterates over all the routes received as a []*Route
// and returns a new *mux.Router ready to use.
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
