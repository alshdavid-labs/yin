package yin

import (
	"net"
	"net/http"
	"strings"
)

func ClientIP(r *http.Request) string {
	clientIP := r.Header.Get(Headers.XForwardedFor)
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get(Headers.XRealIP))
	}
	if clientIP != "" {
		return clientIP
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
