package middleware

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"regexp"
)

const (
	DOMAIN_HEADER = "X-Domain"
	FQDN_HEADER = "X-FQDN"
)

// Sets the X-Domain to be the domain of the callers
// Extracts the domain from the origin header
var Domain = &DomainType{}

type DomainType struct{}

func (a DomainType) Compile(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domain, fqdn := extractDomain(r)
		if r.Header.Get("Origin") == domain {
			log.Errorf("Domain middleware: Unexpected origin %s", domain)
		}
		r.Header.Set(DOMAIN_HEADER, domain)
		r.Header.Set(FQDN_HEADER, fqdn)
		next.ServeHTTP(w, r)
	})
}

func extractDomain(r *http.Request) (domain string, fqdn string) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return origin, origin
	}
	re := regexp.MustCompile(`^https?:\/\/([^:]*)`)
	matches := re.FindAllStringSubmatch(origin, -1)
	if (matches == nil) || (len(matches) != 1) {
		return origin, origin
	}
	groupMatches := matches[0]
	if len(groupMatches) != 2 {
		return origin, origin
	}
	fullDomain := groupMatches[1] //the first match is the whole string, following are the group matches
	re = regexp.MustCompile(`\.`)
	domainParts := re.Split(fullDomain, -1)
	len := len(domainParts)
	switch len {
	case 1:
		domain = fullDomain
	default:
		domain = fmt.Sprintf("%s.%s", domainParts[len-2], domainParts[len-1])
	}
	return domain, fullDomain
}
