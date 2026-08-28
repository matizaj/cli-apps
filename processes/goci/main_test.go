package main

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		out    string
		expErr error
	}{
		{name: "success", proj: "./testdata/tool/", out: "Go Build: SUCCESS\nGo Test: SUCCESS\n", expErr: nil},
		{name: "fail", proj: "./testdata/toolErr/", out: "failed", expErr: &stepErr{step: "go build"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			err := run(tc.proj, &buffer)

			fmt.Printf("tc: %s, output: %v\n", tc.name, buffer.String())

			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected err: %q but got 'nil'", tc.expErr)
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected err: %q, but got %q", tc.expErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected err: %q", err)
			}

			if buffer.String() != tc.out {
				t.Errorf("expected output %s, but got %s", tc.out, buffer.String())
			}
		})
	}
}
