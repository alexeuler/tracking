package middleware

import (
	"github.com/up-finder/silk.web/app/utils/testutils"
	"net/http"
	"testing"
)

func TestCors(t *testing.T) {
	cases := []struct {
		method            string
		header            http.Header
		expectedAllow     string
		expectedServeCall bool
	}{
		{
			method: "POST",
			header: http.Header{
				"Access-Control-Request-Headers": {"content-type, host"},
			},
			expectedAllow:     "content-type, host",
			expectedServeCall: true,
		},
		{
			method: "OPTIONS",
			header: http.Header{
				"Access-Control-Request-Headers": {"content-type, host"},
			},
			expectedAllow:     "content-type, host",
			expectedServeCall: false,
		},
		{
			method: "",
			header: http.Header{
				"Access-Control-Request-Headers": {"content-type"},
			},
			expectedAllow:     "content-type",
			expectedServeCall: true,
		},
		{
			method:            "GET",
			header:            http.Header{},
			expectedAllow:     "",
			expectedServeCall: true,
		},
		{
			method:            "",
			header:            http.Header{},
			expectedAllow:     "",
			expectedServeCall: true,
		},
	}
	for _, c := range cases {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !c.expectedServeCall {
				t.Errorf("Cors Middleware: unexpected call to the next middleware with "+
					"method: %s, header: %v", c.method, c.header)
			}
		})
		corsHandler := Cors.Compile(next)
		r := &http.Request{Header: c.header, Method: c.method}
		r.Header.Set("Origin", "localhost")
		w := testutils.NewFakeResponse()
		corsHandler.ServeHTTP(w, r)
		if w.Header().Get("Access-Control-Allow-Headers") != c.expectedAllow {
			t.Errorf("Cors Middleware: method: (%s), header: (%v): expected Access-Control-Allow-Headers to be (%v) "+
				"got: (%v)", c.method, c.header, c.expectedAllow, w.Header().Get("Access-Control-Allow-Headers"))
		}
		if w.Header().Get("Access-Control-Allow-Origin") != r.Header.Get("Origin") {
			t.Errorf("Cors Middleware: method: (%s), header: (%v): expected Access-Control-Allow-Origin to be (%v) "+
				"got: (%v)", c.method, c.header, r.Header.Get("Origin"), w.Header().Get("Access-Control-Allow-Origin"))
		}
		if w.Header().Get("Access-Control-Allow-Methods") != "POST, GET, OPTIONS" {
			t.Errorf("Cors Middleware: method: (%s), header: (%v): expected Access-Control-Allow-Methods to be (POST, GET, OPTIONS) "+
				"got: (%v)", c.method, c.header, w.Header().Get("Access-Control-Allow-Methods"))
		}
	}
}
