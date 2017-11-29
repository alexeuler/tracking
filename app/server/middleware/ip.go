package middleware

import (
	"net/http"
	"strings"
)

const (
	IP_HEADER = "X-Real-Ip"
)

// Sets the X-Real-Ip header to be the domain of the caller
// Extracts the info from X-Forwarded-For (1st priority),
// X-Real-Ip (2nd priority), r.RemoteAddr (3rd priority)
var Ip = &IpType{}

type IpType struct{}

func (a IpType) Compile(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(IP_HEADER, extractIp(r))
		next.ServeHTTP(w, r)
	})
}

func extractIpWithPort(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ", ")
		return ips[0]
	}
	xri := r.Header.Get("X-Real-Ip")
	if xri != "" {
		return xri
	}
	return r.RemoteAddr
}

func extractIp(r *http.Request) string {
	ip := extractIpWithPort(r)
	idx := strings.LastIndex(ip, ":")
	if idx == -1 {
		return ip
	}
	return ip[:idx]
}
