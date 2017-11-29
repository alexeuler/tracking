package middleware

import (
	"bufio"
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/up-finder/silk.web/app/utils/testutils"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLog(t *testing.T) {
	logger := log.New()
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	logger.Out = writer
	log := &LogType{logger: logger}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	l := log.Compile(next)
	w := testutils.NewFakeResponse()
	body := bytes.NewReader([]byte{})
	r := &http.Request{Body: ioutil.NopCloser(body)}
	l.ServeHTTP(w, r)
	if writer.Buffered() == 0 {
		t.Errorf("Middleware Log: The request was not logged")
	}
}
