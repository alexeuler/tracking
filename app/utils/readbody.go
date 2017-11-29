package utils

import (
	"bytes"
	"io"
	"io/ioutil"
)

// This method is used to read the request body and close it, replacing it with a stub with data
// Useful in the middleware
func ReadBody(body *io.ReadCloser) string {
	bts, _ := ioutil.ReadAll(*body)
	(*body).Close()
	buf := bytes.NewReader(bts)
	*body = ioutil.NopCloser(buf)
	return string(bts)
}
