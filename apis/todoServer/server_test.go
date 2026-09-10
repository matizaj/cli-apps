package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				body []byte
				err error
			)
			resp, err := http.Get(url+tc.path)
			if err != nil {
				t.Error(err)
			}

			defer resp.Body.Close()

			if tc.expCode != resp.StatusCode {
				t.Fatalf("expected %d but got %d", tc.expCode, resp.StatusCode)
			}

			switch {
			case strings.Contains(resp.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(resp.Body); err != nil {
					t.Error(err)
				}
				if !strings.Contains(string(body), tc.expContent) {
					t.Errorf("expected body %s but got %s", tc.expContent, string(body))
				} 
				default:
					t.Errorf("unsupported content-type")
			}
		})
	}
}