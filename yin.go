package yin

import "net/http"

// H is a convenience type
type H map[string]interface{}

func Event(w http.ResponseWriter, r *http.Request) (*Response, *Request) {
	return Res(w, r), Req(r)
}
