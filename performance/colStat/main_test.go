package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		col    int
		op     string
		exp    string
		files  []string
		expErr error
	}{
		{"RunAvg", 1, "a", "1.5", []string{".\\testdata\\example.csv", ".\\testdata\\example2.csv"}, nil},
		{"RunSum", 1, "s", "6", []string{".\\testdata\\example.csv", ".\\testdata\\example2.csv"}, nil},
		{"RunFailedRead", 1, "s", "6", []string{".\\testdata\\example.csv", ".\\testdata\\exampleNotExist.csv"}, os.ErrNotExist},
		{"RunColumntNotExist", 100, "s", "6", []string{".\\testdata\\example.csv"}, ErrInvalidColumn},
		{"RunFailOperation", 1, "sss", "6", []string{".\\testdata\\example.csv"}, ErrInvalidOperation},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			err := run(tc.files, tc.op, tc.col, &buffer)
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

func BenchmarkRun(b *testing.B) {
	filenames, err := filepath.Glob(".\\testdata\\benchmark\\*.csv")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := run(filenames, "a", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}
