package ctest

import (
	"net/http"
	"net/http/httptest"
)

// Request records the http request.
func Request(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
