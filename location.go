package yin

import (
	"net/http"
	"strings"
)

type Location struct {
	Scheme string
	Host   string
	Origin string
}

func getLocation(r *http.Request) *Location {
	scheme := resolveScheme(r)
	host := resolveHost(r)
	origin := scheme + "://" + host
	return &Location{
		Scheme: scheme,
		Host:   host,
		Origin: origin,
	}
}

func resolveScheme(r *http.Request) string {
	customSchemeHeader := r.Header.Get(Headers.XOriginalScheme)
	if customSchemeHeader != "" {
		return customSchemeHeader
	}
	schemeHeader := r.Header.Get(Headers.XForwardedProto)
	if schemeHeader == "https" {
		return schemeHeader
	}
	if r.URL.Scheme == "https" {
		return "https"
	}
	if r.TLS != nil {
		return "https"
	}
	if strings.HasPrefix(r.Proto, "HTTPS") {
		return "https"
	}
	return "http"
}

func resolveHost(r *http.Request) string {
	customHostHeader := r.Header.Get(Headers.XOriginalHost)
	if customHostHeader != "" {
		return customHostHeader
	}
	forwardedForHeader := r.Header.Get(Headers.XForwardedFor)
	if forwardedForHeader != "" {
		return forwardedForHeader
	}
	hostHeader := r.Header.Get("X-Host")
	if hostHeader != "" {
		return hostHeader
	}
	if r.Host != "" {
		return r.Host
	}
	if r.URL.Host != "" {
		return r.URL.Host
	}
	return ""
}
