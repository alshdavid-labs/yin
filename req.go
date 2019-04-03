package yin

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Request struct {
	BindBody    func(body interface{}) error
	GetHeader   func(key string) string
	GetQuery    func(key string) string
	GetCookie   func(name string) string
	GetLocation func() *Location
}

func Req(r *http.Request) *Request {
	req := &Request{}

	req.BindBody = func(body interface{}) error {
		if r.Body == nil {
			return errors.New("No Request body found")
		}
		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			return err
		}
		return nil
	}

	req.GetCookie = func(name string) string {
		cookie, err := r.Cookie(name)
		if err != nil {
			return ""
		}
		val, _ := url.QueryUnescape(cookie.Value)
		return val
	}

	req.GetHeader = func(key string) string {
		return r.Header.Get(key)
	}

	req.GetQuery = func(key string) string {
		return r.URL.Query().Get(key)
	}

	req.GetLocation = func() *Location {
		return getLocation(r)
	}

	return req
}
