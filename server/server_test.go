package server

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tt := []struct {
		name           string
		options        []Option
		expectedToFail bool
	}{
		{
			name:           "with no options",
			expectedToFail: true,
		},
		{
			name: "with BoltDB option",
			options: []Option{
				WithBoltDB("", ""),
			},
		},
		{
			name: "with PostgreSQL option",
			options: []Option{
				WithPostgresDB("", "", "", "", "", ""),
			},
		},
		{
			name: "with BoltDB and PostgreSQL option",
			options: []Option{
				WithBoltDB("", ""),
				WithPostgresDB("", "", "", "", "", ""),
			},
		},
		{
			name: "with testing db",
			options: []Option{
				WithTestDB(),
			},
		},
		{
			name: "with a buggy option",
			options: []Option{
				func(s *server) error {
					return fmt.Errorf("my bad")
				},
			},
			expectedToFail: true,
		},
		{
			name: "with listener and Boltdb",
			options: []Option{
				WithListener(":0"),
				WithBoltDB("", ""),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := New(tc.options...); err != nil {
				if tc.expectedToFail {
					t.Skipf("while creating a new Server failed as expected: %v", err)
				}

				t.Fatalf("while creating a new Server: %v", err)
			}

			if tc.expectedToFail {
				t.Fatalf("while creating a new Server expected to fail creation not failed")
			}
		})
	}
}
