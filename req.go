package yin

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Request struct {
	r *http.Request
}

func Req(r *http.Request) *Request {
	return &Request{r}
}

func (req *Request) BindBody(body interface{}) error {
	if req.r.Body == nil {
		return errors.New("No Request body found")
	}
	err := json.NewDecoder(req.r.Body).Decode(body)
	if err != nil {
		return err
	}
	return nil
}

func (req *Request) GetCookie(name string) string {
	cookie, err := req.r.Cookie(name)
	if err != nil {
		return ""
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val
}

func (req *Request) GetHeader(key string) string {
	return req.r.Header.Get(key)
}

func (req *Request) GetQuery(key string) string {
	return req.r.URL.Query().Get(key)
}

func (req *Request) GetLocation() *Location {
	return getLocation(req.r)
}
