package cache

import "testing"

func TestNew(t *testing.T) {
	tt := []struct {
		name           string
		platform       string
		expectedToFail bool
	}{
		{
			name:     "with platform",
			platform: "ristretto",
		},
		{
			name:           "without platform",
			expectedToFail: true,
		},
	}

	for _, tc := range tt {
		c := Create(tc.platform)
		if c == nil {
			if tc.expectedToFail {
				t.Skipf("nil returned as result as expected")
			}

			t.Fatal("nil value unexpected received")
		}

		if tc.expectedToFail {
			t.Fatalf("test expected to fail did not failed")
		}
	}
}
