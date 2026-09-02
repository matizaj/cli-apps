package scan_test

import (
	"errors"
	"matizaj/cli-apps/cobra/pScan/scan"
	"testing"
)

func TestAdd(t *testing.T) {
	testCases:=[]struct{
		name string
		host string
		expectLen	int
		expectErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hostsList:= scan.HostsList{Hosts: []string{"host1"}}
			err:= hostsList.Add(tc.host)
			if err!= nil {
				if !errors.Is(err, tc.expectErr) {
					t.Fatalf("expected %q but got %q", tc.expectErr, err)
				}
			}

			resultLen := len(hostsList.Hosts)			
			if resultLen != tc.expectLen {
				t.Errorf("expected %d but got %d", tc.expectLen,resultLen)
			}
		})
	}
}