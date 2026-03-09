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

// BuildAnnounceURLForPasskey keeps backward-compatibility and returns the first announce URL.
func BuildAnnounceURLForPasskey(ctx context.Context, passkey string) (string, error) {
	list, err := BuildAnnounceURLsForPasskey(ctx, passkey)
	if err != nil {
		return "", err
	}
	if len(list) == 0 {
		return "", errMissingTrackerHost
	}
	return list[0], nil
}

// BuildAnnounceURLsForPasskey resolves tracker list only from system settings.
func BuildAnnounceURLsForPasskey(ctx context.Context, passkey string) ([]string, error) {
	endpoints := ResolveTrackerEndpointsFromSettings(ctx)
	if len(endpoints) == 0 {
		return nil, errMissingTrackerHost
	}

	return BuildAnnounceURLsFromEndpoints(endpoints, passkey)
}

// BuildAnnounceURLsFromEndpoints builds normalized announce URLs and removes duplicates.
func BuildAnnounceURLsFromEndpoints(endpoints []string, passkey string) ([]string, error) {
	announceURLs := make([]string, 0, len(endpoints))
	seen := map[string]struct{}{}

	for _, endpoint := range endpoints {
		announceURL, err := BuildAnnounceURL(endpoint, passkey)
		if err != nil {
			continue
		}
		if _, ok := seen[announceURL]; ok {
			continue
		}
		seen[announceURL] = struct{}{}
		announceURLs = append(announceURLs, announceURL)
	}

	if len(announceURLs) == 0 {
		return nil, errMissingTrackerHost
	}
	return announceURLs, nil
}

func ResolveTrackerEndpointsFromSettings(ctx context.Context) []string {
	return systemSetting.TrackerListValue(ctx)
}

func parseTrackerList(raw string) []string {
	return systemSetting.ParseTrackerList(raw)
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
