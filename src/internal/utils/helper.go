package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetIpAddressFromRequest(r *http.Request) string {
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		return cfIP
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	forwardedFor := r.Header.Get("X-Forwarded-For")

	if forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return r.RemoteAddr
	}

	return ip
}
func GetHttpHeadersFromRequest(r *http.Request) map[string]string {
	flatHeaders := make(map[string]string)

	for k, v := range r.Header {
		flatHeaders[k] = strings.Join(v, ", ")
	}

	return flatHeaders
}

func GetPublicUrlFromRequest(r *http.Request) string {
	scheme := "http"

	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	return scheme + "://" + r.Host
}
