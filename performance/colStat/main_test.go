package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	testCases:= []struct{
		name string
		col int
		op string
		exp string
		files []string
		expErr error
	}{
		{"RunAvg", 1, "a", "1", []string{".\\testdata\\example.csv"}, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			err := run(tc.files,tc.op, tc.col, &buffer ) 
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected error bu got nil")
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected %q but got %q", tc.expErr, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected err %q", err)
			}

			if buffer.String() != tc.exp {
				t.Errorf("expected %q but got %q", tc.exp, &buffer)
			}
		})
	}
}