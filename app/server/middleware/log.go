package middleware

import (
	log "github.com/Sirupsen/logrus"
	utils "github.com/up-finder/silk.web/app/utils"
	"net/http"
)

// Logs every request
var Log = &LogType{logger: log.StandardLogger()}

type LogType struct {
	logger *log.Logger
}

func (l LogType) Compile(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		url := r.URL
//		host := r.Host
		body := utils.ReadBody(&r.Body)
		l.logger.Infof("Incoming request: %+v with body (%s)", *r, body)
		next.ServeHTTP(w, r)
	})
}
