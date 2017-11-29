package controllers

import (
	"github.com/gorilla/mux"
	"github.com/up-finder/silk.web/app/models"
	"github.com/up-finder/silk.web/app/utils/testutils"
	"net/http"
	"testing"
	"fmt"
	"github.com/up-finder/silk.web/app/server/middleware"
)

func TestScriptShow(t *testing.T) {
	testutils.Setup()
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/script/").Handler(http.HandlerFunc(scriptShow))

	cases := []struct {
		domain   string
		expected string
		code      int
	}{
		{
			domain:   "localhost",
			expected: `script`,
			code: 0,
		},
		{
			domain:   "",
			expected: "",
			code: 500,
		},
	}

	for _, c := range cases {
		r, _ := http.NewRequest("GET", "http://localhost:3000/script/", nil)
		r.Header.Set(middleware.DOMAIN_HEADER, c.domain)
		w := testutils.NewFakeResponse()
		saved := models.FetchScript
		models.FetchScript = func(domain string) (*models.Script, error) {
			if domain == "" {
				return &models.Script{Domain:domain ,Value:""}, fmt.Errorf("")
			} else {
				return &models.Script{Domain:domain, Value:"script"}, nil
			}
		}
		router.ServeHTTP(w, r)
		models.FetchScript = saved
		if (c.code == 0) {
			if w.Body != c.expected {
				t.Errorf("Script Controller: Show: expected response %s, got %s", c.expected, w.Body)
			}
		}
		if c.code!=w.Code {
			t.Errorf("Script Controller: Show: expected code %d, got %d", c.code, w.Code)
		}

	}
}
