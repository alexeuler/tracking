package middleware

import "net/http"

// Every middleware must implement this interface
// Compile function must return handlerfunc
// The idea is to inject next handler into the handlerfunc as a closure and
// implement inside the handler control flow based on call of return or next.ServeHttp
type Middleware interface {
	Compile(next http.Handler) http.Handler
}

// Middleware stack - lastHandler is the exit handler of middleware
type Stack struct {
	data        []Middleware
	lastHandler http.Handler
}

// Creates a new middleware stack
func NewStack(lastHandler http.Handler) *Stack {
	return &Stack{data: make([]Middleware, 0), lastHandler: lastHandler}
}

// Adds middleware to the stack
// The middleware is called according to the order of registration - first registered, first called
func (m *Stack) Register(h ...Middleware) {
	m.data = append(m.data, h...)
}

// Complies the whole stack and returns the handlerfunc, that will execute the whole stack
func (m *Stack) Compile() http.Handler {
	result := m.lastHandler
	for i := len(m.data) - 1; i >= 0; i-- {
		compiler := m.data[i]
		result = compiler.Compile(result)
	}
	return result
}
