package handlers

import "strings"

// RedactPasskeyInPath replaces passkey path segments for safe logging and metrics labels.
// Matches /announce/:passkey, /scrape/:passkey, and /api/announce|scrape/:passkey.
func RedactPasskeyInPath(path string) string {
	if path == "" {
		return path
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := 0; i < len(parts)-1; i++ {
		switch parts[i] {
		case "announce", "scrape":
			parts[i+1] = "[passkey]"
			return "/" + strings.Join(parts, "/")
		}
	}
	return path
}
