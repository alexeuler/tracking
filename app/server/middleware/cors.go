package middleware

import (
	"net/http"
)

// Sets the headers - allow request from any domain, POST, GET, OPTIONS methods, any headers
// Returns if the request was pre-flight (method OPTIONS)
var Cors = &CorsType{}

type CorsType struct{}

func (a CorsType) Compile(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		if r.Method != "OPTIONS" {
			next.ServeHTTP(w, r)
		}
	})
}
