package middleware

import (
	"github.com/up-finder/silk.web/app/db"
	"github.com/up-finder/silk.web/app/models"
	"github.com/up-finder/silk.web/app/utils/testutils"
	"net/http"
	"testing"
)

func TestAuth(t *testing.T) {
	testutils.Setup()

	cases := []struct {
		xdomain           string
		xsecret           string
		redisHash         map[string]string
		expectedServeCall bool
		excludeRoutes     []string
		url               string
	}{
		{
			xdomain:           "localhost",
			xsecret:           "test",
			redisHash:         map[string]string{"localhost": "test"},
			expectedServeCall: true,
		},
		{
			xdomain:           "xyz.ru",
			xsecret:           "test123",
			redisHash:         map[string]string{"xyz.ru": "test123"},
			expectedServeCall: true,
		},
		{
			xdomain:           "",
			xsecret:           "",
			redisHash:         map[string]string{},
			expectedServeCall: false,
		},
		{
			xdomain:           "localhost",
			xsecret:           "",
			redisHash:         map[string]string{},
			expectedServeCall: false,
		},
		{
			xdomain:           "",
			xsecret:           "test123",
			redisHash:         map[string]string{},
			expectedServeCall: false,
		},
		{
			xdomain:           "xyz.ru",
			xsecret:           "test123",
			redisHash:         map[string]string{"xyz.ru": "pass"},
			expectedServeCall: false,
		},
		{
			xdomain:           "",
			xsecret:           "test123",
			redisHash:         map[string]string{},
			expectedServeCall: true,
			excludeRoutes:[]string{"/test"},
			url:"http://localost.ru/test",
		},
		{
			xdomain:           "",
			xsecret:           "test123",
			redisHash:         map[string]string{},
			expectedServeCall: false,
			excludeRoutes:[]string{"/test1"},
			url:"http://localost.ru/test",
		},
	}

	for _, c := range cases {
		db.Redis.FlushDb()
		for k, v := range c.redisHash {
			db.Redis.HSet(models.AUTH_HASH_NAME, k, v)
		}
		exRoutes := c.excludeRoutes
		if (exRoutes == nil) {
			exRoutes = make([]string, 0)
		}
		Auth = &AuthType{excludeRoutes:exRoutes}
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !c.expectedServeCall {
				t.Errorf("Auth Middleware: unexpected call to the next middleware with " +
				"X-Secret: %s, X-Domain: %s, Redis: %v", c.xdomain, c.xsecret, c.redisHash)
			}
		})
		authHandler := Auth.Compile(next)
		w := testutils.NewFakeResponse()
		r, _ := http.NewRequest("GET", c.url, nil);
		r.Header.Add("X-Secret", c.xsecret);
		r.Header.Add("X-Domain", c.xdomain);
		authHandler.ServeHTTP(w, r)
		if !c.expectedServeCall && w.Code != http.StatusUnauthorized {
			t.Errorf("Auth Middleware: Expected unauthorized, but got: %v for X-Secret: %s, X-Domain: %s, Redis: %v",
				w.Code, c.xdomain, c.xsecret, c.redisHash)
		}
	}
}
