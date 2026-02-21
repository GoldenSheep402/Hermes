package handlers

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP extracts the real client IP from the request, preferring X-Forwarded-For.
// It gracefully falls back to RemoteAddr when no proxy headers are present.
func GetClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// XFF may contain multiple comma-separated values; take the first one.
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		// Strip optional port if present.
		host, _, err := net.SplitHostPort(ip)
		if err == nil {
			return host
		}
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
