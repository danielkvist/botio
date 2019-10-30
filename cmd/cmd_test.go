package cmd

import "testing"

func TestCheckURL(t *testing.T) {
	tt := []struct {
		name           string
		originalURL    string
		expectedURL    string
		expectedToFail bool
	}{
		{
			name:           "base URL",
			originalURL:    "google.com",
			expectedURL:    "http://google.com/api/commands",
			expectedToFail: false,
		},
		{
			name:           "no URL",
			expectedToFail: true,
		},
		{
			name:           "http URL",
			originalURL:    "http://google.com",
			expectedURL:    "http://google.com/api/commands",
			expectedToFail: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			u, err := checkURL(tc.originalURL, true, true)
			if err != nil {
				if tc.expectedToFail {
					t.Logf("while checking URL failed as expected: %v", err)
					return
				}
				t.Fatalf("while checking URL not expected to fail failed: %v", err)
			}

			if u != tc.expectedURL {
				t.Fatalf("expected checked URL to be %q. got=%q", tc.expectedURL, u)
			}
		})
	}
}
