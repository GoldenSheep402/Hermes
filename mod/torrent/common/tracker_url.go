package common

import (
	"context"
	"errors"
	"net/url"
	"strings"

	systemSetting "github.com/GoldenSheep402/Hermes/mod/system/setting"
)

var (
	errMissingPasskey     = errors.New("missing passkey")
	errMissingTrackerHost = errors.New("missing tracker host")
)

// BuildAnnounceURL builds an announce URL with the given passkey.
// endpoint can be a host, a base URL, or an announce URL.
func BuildAnnounceURL(endpoint string, passkey string) (string, error) {
	passkey = strings.TrimSpace(passkey)
	if passkey == "" {
		return "", errMissingPasskey
	}

	raw := strings.TrimSpace(endpoint)
	if raw == "" {
		return "", errMissingTrackerHost
	}

	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(u.Host) == "" {
		return "", errMissingTrackerHost
	}

	path := normalizeAnnouncePath(u.Path)
	u.Path = path + "/" + passkey
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

// BuildAnnounceURLForPasskey resolves the single site announce URL from system settings.
func BuildAnnounceURLForPasskey(ctx context.Context, passkey string) (string, error) {
	endpoint := ResolveTrackerEndpointFromSettings(ctx)
	if endpoint == "" {
		return "", errMissingTrackerHost
	}
	return BuildAnnounceURL(endpoint, passkey)
}

// BuildAnnounceURLsForPasskey returns a one-element list for callers that still expect a slice.
func BuildAnnounceURLsForPasskey(ctx context.Context, passkey string) ([]string, error) {
	announceURL, err := BuildAnnounceURLForPasskey(ctx, passkey)
	if err != nil {
		return nil, err
	}
	return []string{announceURL}, nil
}

// ResolveTrackerEndpointFromSettings returns the configured public announce endpoint.
func ResolveTrackerEndpointFromSettings(ctx context.Context) string {
	return systemSetting.TrackerAnnounceURLValue(ctx)
}

func normalizeAnnouncePath(path string) string {
	p := strings.TrimSpace(path)
	p = strings.TrimSuffix(p, "/")
	if p == "" {
		return "/announce"
	}

	if idx := strings.Index(p, "/announce/"); idx >= 0 {
		return strings.TrimSuffix(p[:idx+len("/announce")], "/")
	}
	if strings.HasSuffix(p, "/announce") {
		return p
	}
	return p + "/announce"
}
