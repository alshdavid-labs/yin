package yin

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	w http.ResponseWriter
	r *http.Request
}

func Res(w http.ResponseWriter, r *http.Request) *Response {
	return &Response{w, r}
}

func (res *Response) SetStatus(statusCode int) *Response {
	res.w.WriteHeader(statusCode)
	return res
}

func (res *Response) SetCookie(cookie *http.Cookie) *Response {
	http.SetCookie(res.w, cookie)
	return res
}

func (res *Response) SetHeader(key string, value string) *Response {
	res.w.Header().Set(key, value)
	return res
}

func (res *Response) SendJSON(u interface{}) {
	res.w.Header().Set(Headers.ContentType, "application/json")
	json.NewEncoder(res.w).Encode(u)
}

func (res *Response) SendString(s string) {
	res.w.Write([]byte(s))
}

func (res *Response) SendStatus(statusCode int) {
	res.w.WriteHeader(statusCode)
	res.w.Write([]byte(""))
}

func (res *Response) SendFile(filepath string) {
	http.ServeFile(res.w, res.r, filepath)
}

func (res *Response) SendRedirect(statusCode int, url string) {
	http.Redirect(res.w, res.r, url, statusCode)
}
