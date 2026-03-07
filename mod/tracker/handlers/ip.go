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

// IsLANIP checks if the given IP address is a private LAN IP (e.g. 10.0.0.0/8, 192.168.0.0/16, loopback)
// or belongs to any of the allowed subnets specified in the configuration.
func IsLANIP(ipStr string, allowedSubnets []string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check standard private / loopback / link-local unicast IPs
	if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
		return true
	}

	// Check custom allowed subnets
	for _, cidr := range allowedSubnets {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err == nil && ipNet.Contains(ip) {
			return true
		}
	}

	return false
}

// ExtractIPs retrieves the real public IP from the request and a valid LAN IP
// if provided in the query string ("ip" or "ipv4").
func ExtractIPs(req *http.Request, allowedSubnets []string) (string, string) {
	realIP := GetClientIP(req)
	lanIP := ""

	queryIP := req.URL.Query().Get("ip")
	if queryIP == "" {
		queryIP = req.URL.Query().Get("ipv4")
	}

	if queryIP != "" && IsLANIP(queryIP, allowedSubnets) {
		lanIP = queryIP
	}

	return realIP, lanIP
}
