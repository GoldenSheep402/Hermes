package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientIP_FromXForwardedFor(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err)
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 10.0.0.1")

	ip := GetClientIP(req)
	assert.Equal(t, "203.0.113.1", ip)
}

func TestGetClientIP_FromRemoteAddr(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err)
	req.RemoteAddr = "192.0.2.5:12345"

	ip := GetClientIP(req)
	assert.Equal(t, "192.0.2.5", ip)
}

func TestIsLANIP(t *testing.T) {
	tests := []struct {
		name           string
		ip             string
		allowedSubnets []string
		expected       bool
	}{
		{"Standard Private IP 192", "192.168.1.100", nil, true},
		{"Standard Private IP 10", "10.0.0.1", nil, true},
		{"Standard Private IP 172", "172.16.0.1", nil, true},
		{"Public IP", "8.8.8.8", nil, false},
		{"Allowed custom subnet", "111.222.33.44", []string{"111.222.33.0/24"}, true},
		{"Not in custom subnet", "111.222.34.44", []string{"111.222.33.0/24"}, false},
		{"Invalid IP", "invalid-ip", nil, false},
		{"Loopback IP", "127.0.0.1", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsLANIP(tt.ip, tt.allowedSubnets))
		})
	}
}

func TestExtractIPs(t *testing.T) {
	tests := []struct {
		name           string
		remoteAddr     string
		xForwardedFor  string
		queryIP        string
		allowedSubnets []string
		expectedReal   string
		expectedLan    string
	}{
		{
			name:         "Real IP only",
			remoteAddr:   "200.200.200.1:1234",
			expectedReal: "200.200.200.1",
			expectedLan:  "",
		},
		{
			name:          "Real IP + XFF",
			remoteAddr:    "10.0.0.1:1234",
			xForwardedFor: "200.200.200.1",
			expectedReal:  "200.200.200.1",
			expectedLan:   "",
		},
		{
			name:         "Real IP + Valid LAN IP via query",
			remoteAddr:   "200.200.200.1:1234",
			queryIP:      "192.168.1.10",
			expectedReal: "200.200.200.1",
			expectedLan:  "192.168.1.10",
		},
		{
			name:         "Real IP + Invalid LAN IP (public IP spoofing)",
			remoteAddr:   "200.200.200.1:1234",
			queryIP:      "8.8.8.8",
			expectedReal: "200.200.200.1",
			expectedLan:  "", // should be rejected!
		},
		{
			name:           "Real IP + Custom Allowed Subnet",
			remoteAddr:     "200.200.200.1:1234",
			queryIP:        "111.222.33.44",
			allowedSubnets: []string{"111.222.33.0/24"},
			expectedReal:   "200.200.200.1",
			expectedLan:    "111.222.33.44",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "http://example.com/announce"
			if tt.queryIP != "" {
				url += "?ip=" + tt.queryIP
			}
			req, _ := http.NewRequest("GET", url, nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			realIP, lanIP := ExtractIPs(req, tt.allowedSubnets)
			assert.Equal(t, tt.expectedReal, realIP, "Real IP mismatch")
			assert.Equal(t, tt.expectedLan, lanIP, "Lan IP mismatch")
		})
	}
}
