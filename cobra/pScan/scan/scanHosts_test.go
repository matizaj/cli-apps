package scan_test

import (
	"matizaj/cli-apps/cobra/pScan/scan"
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
