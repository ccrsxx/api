package utils

import (
	"net/http"
	"strings"
)

func GetIpAddressFromRequest(r *http.Request) string {
	ipAddress := r.Header.Get("X-Forwarded-For")

	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	return ipAddress
}

func GetHttpHeadersFromRequest(r *http.Request) map[string]string {
	flatHeaders := make(map[string]string)

	for k, v := range r.Header {
		flatHeaders[k] = strings.Join(v, ", ")
	}

	return flatHeaders
}
