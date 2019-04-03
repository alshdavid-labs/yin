package yin

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	SetStatus    func(statusCode int) *Response
	SetCookie    func(cookie *http.Cookie) *Response
	SetHeader    func(key string, value string) *Response
	SendJSON     func(interface{})
	SendString   func(s string)
	SendStatus   func(statusCode int)
	SendFile     func(filepath string)
	SendRedirect func(statusCode int, url string)
}

func Res(w http.ResponseWriter, r *http.Request) *Response {
	res := &Response{}

	res.SetStatus = func(statusCode int) *Response {
		w.WriteHeader(statusCode)
		return res
	}

	res.SetCookie = func(cookie *http.Cookie) *Response {
		http.SetCookie(w, cookie)
		return res
	}

	res.SetHeader = func(key string, value string) *Response {
		w.Header().Set(key, value)
		return res
	}

	res.SendJSON = func(u interface{}) {
		w.Header().Set(Headers.ContentType, "application/json")
		json.NewEncoder(w).Encode(u)
	}

	res.SendString = func(s string) {
		w.Write([]byte(s))
	}

	res.SendStatus = func(statusCode int) {
		w.WriteHeader(statusCode)
		w.Write([]byte(""))
	}

	res.SendFile = func(filepath string) {
		http.ServeFile(w, r, filepath)
	}

	res.SendRedirect = func(statusCode int, url string) {
		http.Redirect(w, r, url, statusCode)
	}

	return res
}
