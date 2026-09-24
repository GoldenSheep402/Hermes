package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildAnnounceURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		endpoint string
		passkey  string
		expected string
	}{
		{
			name:     "host only",
			endpoint: "tracker.hermes.local:8080",
			passkey:  "pk123",
			expected: "http://tracker.hermes.local:8080/announce/pk123",
		},
		{
			name:     "base path",
			endpoint: "https://tracker.hermes.local/tracker",
			passkey:  "pk123",
			expected: "https://tracker.hermes.local/tracker/announce/pk123",
		},
		{
			name:     "announce path",
			endpoint: "https://tracker.hermes.local/announce",
			passkey:  "pk123",
			expected: "https://tracker.hermes.local/announce/pk123",
		},
		{
			name:     "announce with old passkey replaced",
			endpoint: "https://tracker.hermes.local/announce/oldpk",
			passkey:  "newpk",
			expected: "https://tracker.hermes.local/announce/newpk",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual, err := BuildAnnounceURL(tc.endpoint, tc.passkey)
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestBuildAnnounceURLValidation(t *testing.T) {
	t.Parallel()

	_, err := BuildAnnounceURL("", "pk")
	require.Error(t, err)

	_, err = BuildAnnounceURL("tracker.hermes.local", "")
	require.Error(t, err)
}
