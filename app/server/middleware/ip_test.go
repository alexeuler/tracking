package middleware

import (
	"github.com/up-finder/silk.web/app/utils/testutils"
	"net/http"
	"testing"
)

func TestIp(t *testing.T) {
	cases := []struct {
		xRealIp       string
		xForwarderFor string
		remote        string
		expected      string
	}{
		{
			xRealIp:  "127.0.0.1:5000",
			expected: "127.0.0.1",
		},
		{
			xRealIp:       "127.0.0.1:5000",
			xForwarderFor: "192.186.0.1:8000, 127.0.0.1:5000",
			expected:      "192.186.0.1",
		},
		{
			xRealIp:       "127.0.0.1",
			xForwarderFor: "192.186.0.1, 127.0.0.1",
			remote:        "192.127.1.1",
			expected:      "192.186.0.1",
		},
		{
			remote:   "192.127.1.1:8000",
			expected: "192.127.1.1",
		},
		{
			remote:   "[::1]:8000",
			expected: "[::1]",
		},
	}

	for _, c := range cases {
		r := &http.Request{Header: http.Header{}}
		if c.xRealIp != "" {
			r.Header.Set("X-Real-Ip", c.xRealIp)
		}
		if c.xForwarderFor != "" {
			r.Header.Set("X-Forwarded-For", c.xForwarderFor)
		}
		r.RemoteAddr = c.remote
		w := testutils.NewFakeResponse()
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
		ipHandler := Ip.Compile(next)
		ipHandler.ServeHTTP(w, r)
		if r.Header.Get(IP_HEADER) != c.expected {
			t.Errorf("Domain Middleware: case (%+v): expected (%s), got (%s)",
				c, c.expected, r.Header.Get("X-Real-Ip"))
		}
	}
}
