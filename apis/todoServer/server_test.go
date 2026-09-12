package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setupApi(t *testing.T) (string, func()) {
	t.Helper()
	tempTodoFile, err := os.CreateTemp("", "todoTest")
	if err != nil {
		t.Fatal(err)
	}
	defer tempTodoFile.Close()

	ts:=httptest.NewServer(newMux(tempTodoFile.Name()))

	for i:=1 ; i<3 ; i++ {
		var body bytes.Buffer
		taskName := fmt.Sprintf("Task Number - %d", i)
		item := struct{
			Task string `json:task`
		}{
			Task: taskName,
		}

		if err := json.NewEncoder(&body).Encode(item); err!= nil {
			t.Fatal(err)
		}

		r, err := http.Post(ts.URL+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}

		if r.StatusCode != http.StatusCreated {
			t.Fatalf("failed to add initial data %d", r.StatusCode)
		}
	}
	return ts.URL, func() {
		ts.Close()
		os.Remove(tempTodoFile.Name())
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