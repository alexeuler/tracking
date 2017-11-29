package server

import (
	"github.com/gorilla/mux"
	"github.com/up-finder/silk.web/app/controllers"
	"net/http"
)

// This is 3rd party gorilla router, used as the last entry in the middleware stack
type Route struct {
	Method      string           //http method
	Pattern     string           //pattern for matching
	HandlerFunc http.HandlerFunc //handler
	Name        string           //name for the named route
}

type Routes []Route

var routes = Routes{
	Route{
		"GET",
		"/events/{id}",
		controllers.Event.Show,
		"EventShow",
	},
	Route{
		"POST",
		"/events",
		controllers.Event.Create,
		"EventCreate",
	},
	Route{
		"GET",
		"/script",
		controllers.Script.Show,
		"ScriptShow",
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Handler(handler).
			Name(route.Name)
	}
	return router
}
