package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/tools/go/analysis/passes/defers"
)

func setupApi(t *testing.T) (string, func()) {
	t.Helper()

	ts:=httptest.NewServer(newMux(""))
	return ts.URL, func() {
		ts.Close()
	}
}

func TestGet(t *testing.T) {
	testCases := []struct{
		name string
		path string
		expCode int
		expItems int
		expContent string
	}{
		{"GetRoot", "/", http.StatusOK, 0, "todo server api"},
		{"NotFound", "/todo/500", http.StatusNotFound, 0, ""},
	}

	url, cleanup := setupApi(t)
	defer cleanup()
}