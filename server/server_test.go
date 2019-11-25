package server

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	tt := []struct {
		name           string
		options        []Option
		expectedToFail bool
		errMsg         string
	}{
		{
			name: "with basic options",
			options: []Option{
				WithTestDB(),
				WithRistrettoCache(1 << 30),
				WithHTTPPort(":8081"),
				WithListener(":0"),
				WithInsecureGRPCServer(),
				WithTextLogger(&bytes.Buffer{}),
			},
		},
		{
			name: "with repeated options",
			options: []Option{
				WithTestDB(),
				WithTestDB(),
				WithRistrettoCache(1 << 30),
				WithRistrettoCache(200),
				WithHTTPPort(":8081"),
				WithListener(":0"),
				WithInsecureGRPCServer(),
				WithTextLogger(&bytes.Buffer{}),
				WithJSONLogger(&bytes.Buffer{}),
			},
		},
		{
			name:           "zero options",
			options:        []Option{},
			expectedToFail: true,
			errMsg:         "while creating a new Server: no options provided",
		},
		{
			name: "without DB",
			options: []Option{
				WithRistrettoCache(1 << 30),
				WithHTTPPort(":8081"),
				WithListener(":0"),
				WithInsecureGRPCServer(),
				WithTextLogger(&bytes.Buffer{}),
			},
			expectedToFail: true,
			errMsg:         "while creating a new Server: no DB provided",
		},
		{
			name: "without Cache",
			options: []Option{
				WithTestDB(),
				WithHTTPPort(":8081"),
				WithListener(":0"),
				WithInsecureGRPCServer(),
				WithTextLogger(&bytes.Buffer{}),
			},
			expectedToFail: true,
			errMsg:         "while creating a new Server: no Cache provided",
		},
		{
			name: "without Listener",
			options: []Option{
				WithTestDB(),
				WithRistrettoCache(1 << 30),
				WithHTTPPort(":8081"),
				WithInsecureGRPCServer(),
				WithTextLogger(&bytes.Buffer{}),
			},
			expectedToFail: true,
			errMsg:         "while creating a new Server: no net.Listener provided",
		},
		{
			name: "without gRPC server",
			options: []Option{
				WithTestDB(),
				WithRistrettoCache(1 << 30),
				WithHTTPPort(":8081"),
				WithListener(":0"),
				WithTextLogger(&bytes.Buffer{}),
			},
			expectedToFail: true,
			errMsg:         "while creating a new Server: no gRPC server provided",
		},
		{
			name: "without logger",
			options: []Option{
				WithTestDB(),
				WithRistrettoCache(1 << 30),
				WithHTTPPort(":8081"),
				WithListener(":0"),
				WithInsecureGRPCServer(),
			},
			expectedToFail: true,
			errMsg:         "while creating a new Server: no logger provided",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(tc.options...)
			if err != nil {
				if tc.expectedToFail && tc.errMsg == err.Error() {
					t.Skipf("test failed as expected: %v", err)
				}

				t.Fatalf("while creating a new Server: %v", err)
			}

			if tc.expectedToFail {
				t.Fatalf("test did not failed as was expected with error: %s", tc.errMsg)
			}
		})
	}
}
