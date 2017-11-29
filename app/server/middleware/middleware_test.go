package middleware

import (
	"github.com/up-finder/silk.web/app/utils/testutils"
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	tracker := []string{}
	first := MiddlewareStub(&tracker, "first")
	second := MiddlewareStub(&tracker, "second")
	initial := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracker = append(tracker, "initial")
	})
	stack := NewStack(initial)
	stack.Register(first, second)
	w := testutils.NewFakeResponse()
	stack.Compile().ServeHTTP(w, nil)
	if (tracker[0] != "first") || (tracker[1] != "second") || (tracker[2] != "initial") {
		t.Errorf("Test Middleware: expected order (first, second, initial), got (%v)", tracker)
	}
}

type MWStub struct {
	tracker *[]string
	message string
}

func (m MWStub) Compile(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*m.tracker = append(*m.tracker, m.message)
		next.ServeHTTP(w, r)
	})
}

func MiddlewareStub(tracker *[]string, message string) MWStub {
	return MWStub{tracker: tracker, message: message}
}
