package middleware

import (
	"github.com/up-finder/silk.web/app/utils/testutils"
	"net/http"
	"testing"
)

func TestDomain(t *testing.T) {
	cases := []struct {
		origin   string
		expectedDomain string
		expectedFQDN string
	}{
		{
			origin:   "http://xyz.ru",
			expectedDomain: "xyz.ru",
			expectedFQDN: "xyz.ru",
		},
		{
			origin:   "http://abc.xyz.ru",
			expectedDomain: "xyz.ru",
			expectedFQDN: "abc.xyz.ru",
		},
		{
			origin:   "https://abc.sfd.xyz.ru:8000",
			expectedDomain: "xyz.ru",
			expectedFQDN: "abc.sfd.xyz.ru",
		},
		{
			origin:   "chrome://xyz.ru",
			expectedDomain: "chrome://xyz.ru",
			expectedFQDN: "chrome://xyz.ru",
		},
		{
			origin:   "",
			expectedDomain: "",
			expectedFQDN: "",
		},
	}

	for _, c := range cases {
		r := &http.Request{Header: http.Header{"Origin": {c.origin}}}
		w := testutils.NewFakeResponse()
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
		domainHandler := Domain.Compile(next)
		domainHandler.ServeHTTP(w, r)
		if r.Header.Get(DOMAIN_HEADER) != c.expectedDomain {
			t.Errorf("Domain Middleware: origin: (%s): expected domain (%s), got (%s)",
				c.origin, c.expectedDomain, r.Header.Get(DOMAIN_HEADER))
		}
		if r.Header.Get(FQDN_HEADER) != c.expectedFQDN {
			t.Errorf("Domain Middleware: origin: (%s): expected fqdn (%s), got (%s)",
				c.origin, c.expectedFQDN, r.Header.Get(FQDN_HEADER))
		}
	}
}
