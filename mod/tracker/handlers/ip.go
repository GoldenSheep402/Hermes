package handlers

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP extracts the client IP from the request.
// When trustedProxyCIDRs is non-empty, X-Forwarded-For is used only if RemoteAddr matches
// a trusted CIDR; otherwise RemoteAddr is used (prevents forged XFF from direct clients).
// When trustedProxyCIDRs is empty, legacy behavior applies: prefer the first XFF hop if present.
func GetClientIP(r *http.Request, trustedProxyCIDRs []string) string {
	remoteIP := parseRequestRemoteIP(r)
	xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))

	if len(trustedProxyCIDRs) == 0 {
		if xff != "" {
			return parseFirstXFF(xff)
		}
		if remoteIP != nil {
			return remoteIP.String()
		}
		return stripHostPortFallback(r.RemoteAddr)
	}

	if remoteIP == nil || !ipMatchesTrustedProxies(remoteIP, trustedProxyCIDRs) {
		if remoteIP != nil {
			return remoteIP.String()
		}
		return stripHostPortFallback(r.RemoteAddr)
	}

	if xff != "" {
		return parseFirstXFF(xff)
	}
	// remoteIP is non-nil here: otherwise the trusted-proxy branch above would have returned.
	return remoteIP.String()
}

func parseRequestRemoteIP(r *http.Request) net.IP {
	host := stripHostPortFallback(r.RemoteAddr)
	if host == "" {
		return nil
	}
	return net.ParseIP(host)
}

func stripHostPortFallback(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err == nil {
		return host
	}
	return addr
}

func parseFirstXFF(xff string) string {
	parts := strings.Split(xff, ",")
	ip := strings.TrimSpace(parts[0])
	host, _, err := net.SplitHostPort(ip)
	if err == nil {
		return host
	}
	return ip
}

func ipMatchesTrustedProxies(ip net.IP, cidrs []string) bool {
	if ip == nil {
		return false
	}
	for _, c := range cidrs {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		_, ipNet, err := net.ParseCIDR(c)
		if err != nil {
			// Allow single IP as /32 or /128
			parsed := net.ParseIP(c)
			if parsed == nil {
				continue
			}
			if parsed.Equal(ip) {
				return true
			}
			continue
		}
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}
