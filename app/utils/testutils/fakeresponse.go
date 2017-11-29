package testutils

import "net/http"

type FakeResponse struct {
	Body   string
	header http.Header
	Code   int
}

func (r *FakeResponse) Header() http.Header {
	return r.header
}

func (r *FakeResponse) Write(data []byte) (int, error) {
	r.Body = string(data)
	return 0, nil
}

func (r *FakeResponse) WriteHeader(x int) {
	r.Code = x
}

func NewFakeResponse() *FakeResponse {
	return &FakeResponse{header: http.Header{}}
}
