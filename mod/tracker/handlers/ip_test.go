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
