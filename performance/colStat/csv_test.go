package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"testing/iotest"
)

func TestMaxValue(t *testing.T) {
	testCases := []struct {
		name string
		data []float64
		exp float64
	}{
		{name: "Max", data: []float64{10, 20, 15, 30, 45, 50, 100, 30}, exp: 100},
		{name: "Max", data: []float64{10, -20, -15, -30, -45, -50, -100, -30}, exp: 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := max(tc.data)

			if res != tc.exp {
				t.Errorf("expected %v, but got %v\n", tc.exp, res)
			}
		})
	}
}
func TestOperations(t *testing.T) {
	data := [][]float64{
		{10, 20, 15, 30, 45, 50, 100, 30},
		{5.5, 8, 2.2, 9.75, 8.45, 3, 2.5, 10.25, 4.75, 6.1, 7.67, 12.287, 5.47},
		{-10, -20},
		{102, 37, 44, 57, 67, 129},
	}

	testCases := []struct {
		name string
		op   statsFunc
		exp  []float64
	}{
		{"sum", sum, []float64{300, 85.927, -30, 436}},
		{"avg", avg, []float64{37.5, 6.609769230769231, -15, 72.666666666666666}},
	}

	for _, tc := range testCases {
		for k, exp := range tc.exp {
			name := fmt.Sprintf("%sData%d", tc.name, k)
			t.Run(name, func(t *testing.T) {
				res := tc.op(data[k])
				if res != exp {
					t.Errorf("expected %g, got %g", exp, res)
				}
			})
		}
	}
}

func TestCsv2Float(t *testing.T) {
	csvData := `IP Address,Requests,Response Time
192.168.0.199,2056,236
192.168.0.88,899,220
192.168.0.199,3054,226
192.168.0.100,4133,218
192.168.0.199,950,238
`
	testCases := []struct {
		name   string
		col    int
		exp    []float64
		expErr error
		r      io.Reader
	}{
		{"Column2", 2, []float64{2056, 899, 3054, 4133, 950}, nil, bytes.NewBufferString(csvData)},
		{"Column3", 3, []float64{236, 220, 226, 218, 238}, nil, bytes.NewBufferString(csvData)},
		{"FailRead", 1, nil, iotest.ErrTimeout, iotest.TimeoutReader(bytes.NewReader([]byte{0}))},
		{"FailedNotNumber", 1, nil, ErrNotNumber, bytes.NewBufferString(csvData)},
		{"FailedInvlidCOlumn", 4, nil, ErrInvalidColumn, bytes.NewBufferString(csvData)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := csv2float(tc.r, tc.col)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected error but got nil instead")
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected %q but got %q", tc.expErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected err %q", err)
			}

			for i, exp := range tc.exp {
				if res[i] != exp {
					t.Errorf("expected %g but got %g", exp, res[i])
				}
			}
		})
	}
}
