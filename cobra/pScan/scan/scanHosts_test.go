package scan_test

import (
	"matizaj/cli-apps/cobra/pScan/scan"
	"net"
	"strconv"
	"testing"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}
	if ps.Open.String() != "closed" {
		t.Errorf("Expected %q, got %q instead​​\n​", "closed", ps.Open.String())
	}

	ps.Open = true
	if ps.Open.String() != "open" {
		t.Errorf("Expected %q, got %q instead​​\n​", "closed", ps.Open.String())
	}
}


func  TestRunHostFound(t *testing.T) {
	testCases:= []struct{
		name string
		expectedState string
	}{
		{"OpenPort","open"},
		{"ClosedPort","closed"},
	}

	host:= "localhost"
	hl:=&scan.HostsList{}
	hl.Add(host)

	ports:=[]int{}
	for _, tc:=range testCases {
		ln, err := net.Listen("tcp", net.JoinHostPort(host,"0"))
		if err!= nil {
			t.Fatal(err)
		}

		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err!= nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err!= nil {
			t.Fatal(err)
		}

		ports=append(ports, port)

		if tc.name == "ClosedPort" {
			ln.Close()
		}

		res := scan.Run(hl, ports)
		if len(res) != 1 {
			t.Fatalf("expected 1 but got %d", len(res))
		}

		if res[0].Host != host {
			t.Errorf("expected host  but got %s", res[0].Host)
		}

		if res[0].NotFound {
			t.Errorf("expected host  to be found %s", res[0].Host)	
		}

		if len(res[0].PortStates) != 2 {
			t.Errorf("expected 2 port states but got %d", len(res[0].PortStates))
		}
	}
}