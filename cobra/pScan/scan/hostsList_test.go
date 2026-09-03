package scan_test

import (
	"errors"
	"matizaj/cli-apps/cobra/pScan/scan"
	"os"
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

func TestDelete(t * testing.T) {
	testCases := []struct{
		name string
		host string
		expectedLen int
		expectedErr error
	}{
		{"RemoveExisting", "host2", 1, nil},
		{"RemoveNotFound", "host3", 2, scan.ErrNotExists},
	}

	for _, tc :=range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hostsList:=scan.HostsList{Hosts: []string{"host1", "host2"}}
			err := hostsList.Remove(tc.host)
			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected error %q but got %q", tc.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error but got %q", err)
			}

			if len(hostsList.Hosts) != tc.expectedLen {
				t.Fatalf("expected list length %d but %d insted", tc.expectedLen, len(hostsList.Hosts))
			}
		})
	}
}

func TestSaveLoad(t *testing.T){
	hl1:=scan.HostsList{}
	hl2:=scan.HostsList{}
	hostname:="host1"
	filename:= "host-list.txt"
	hl1.Add(hostname)

	file, err := os.CreateTemp("", filename)
	if err != nil {
		t.Fatalf("unexpected err %q",err)
	}

	defer file.Close()

	hl1.Save(file.Name())

	err = hl2.Load(file.Name())
	if err != nil {
		t.Fatalf("unexpected err for loading %q", err)
	}

	if len(hl1.Hosts) != len(hl2.Hosts) {
		t.Fatalf("expected %d but got %d",len(hl1.Hosts), len(hl2.Hosts))
	}
}