package middleware

import (
	"github.com/Sirupsen/logrus"
	"github.com/up-finder/silk.web/app/models"
	"net/http"
)

const (
	AUTH_SECRET_FIELD  = "secret"
	AUTH_SECRET_HEADER = "X-Secret"
)

// Matches the domain key passed in X-Secret header with secret stored in models.Auth
// Empty domains and secrets are not allowed
// If unauthorized returns http.StatusUnauthorized, o/w forwards the request further up the stack
var Auth = &AuthType{
	excludeRoutes:[]string{"/script"},
}

type AuthType struct{
	excludeRoutes []string
}

func (a AuthType) Compile(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if (r!=nil) && (r.URL!=nil) {
			path := r.URL.Path
			for _,s:=range a.excludeRoutes {
				if s==path {
					next.ServeHTTP(w,r)
					return
				}
			}
		}

		domain := r.Header.Get(DOMAIN_HEADER)
		auth := models.FetchAuth(domain)
		dbkey := auth.Key

		key := r.Header.Get(AUTH_SECRET_HEADER)

		if (key == "") || (domain == "") || (key != dbkey) {
			logrus.Debugf("Unauthorized: Received key %s for domain %s", key, domain)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
