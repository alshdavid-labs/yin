package yin

import (
	"encoding/json"
	"net/http"
)

type MockHTTPWriter struct {
	Written    []byte
	StatusCode int
}

func (m *MockHTTPWriter) Write(w []byte) (int, error) {
	m.Written = w
	return 0, nil
}

func (m *MockHTTPWriter) Header() http.Header {
	return http.Header{}
}

func (m *MockHTTPWriter) WriteHeader(statusCode int) {
	m.StatusCode = statusCode
}

func (m *MockHTTPWriter) GetBodyJSON() H {
	var v H
	json.Unmarshal(m.Written, &v)
	return v
}

func (m *MockHTTPWriter) GetBodyString() string {
	return string(m.Written)
}
