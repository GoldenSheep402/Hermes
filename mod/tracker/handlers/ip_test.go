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

	ip := GetClientIP(req, nil)
	assert.Equal(t, "203.0.113.1", ip)
}

func TestGetClientIP_FromRemoteAddr(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err)
	req.RemoteAddr = "192.0.2.5:12345"

	ip := GetClientIP(req, nil)
	assert.Equal(t, "192.0.2.5", ip)
}

func TestGetClientIP_TrustedProxyIgnoredWhenDirect(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err)
	req.RemoteAddr = "198.51.100.2:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.99")

	ip := GetClientIP(req, []string{"127.0.0.1/32"})
	assert.Equal(t, "198.51.100.2", ip)
}

func TestGetClientIP_TrustedProxyUsesXFF(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(t, err)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("X-Forwarded-For", "203.0.113.5")

	ip := GetClientIP(req, []string{"127.0.0.1/32"})
	assert.Equal(t, "203.0.113.5", ip)
}
