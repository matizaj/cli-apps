package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"matizaj/cli-apps/todo"
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
		{"GetAll", "/todo", http.StatusOK, 2, "Task Number - 1"},
		{"GetOne", "/todo/1", http.StatusOK, 1, "Task Number - 1"},
		{"NotFound", "/todo/500", http.StatusBadRequest, 0, ""},
	}

	url, cleanup := setupApi(t)
	defer cleanup()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				resp struct {
					Results todo.List `json:results`
					Date int64 `json:date`
					TotalResults int 	`json:total_results`
				}
				body []byte
				err error
			)
			r, err := http.Get(url+tc.path)
			if err != nil {
				t.Error(err)
			}

			defer r.Body.Close()

			if tc.expCode != r.StatusCode {
				t.Fatalf("expected %d but got %d", tc.expCode, r.StatusCode)
			}

			switch {
			case strings.Contains(r.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(r.Body); err != nil {
					t.Error(err)
				}
				if !strings.Contains(string(body), tc.expContent) {
					t.Errorf("expected body %s but got %s", tc.expContent, string(body))
				} 
				case strings.Contains(r.Header.Get("Content-Type"), "application/json"):
					if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
						t.Error(err)
					}
					if resp.TotalResults != tc.expItems {
						t.Errorf("expected %d got %d", tc.expItems, resp.TotalResults)
					}
				default:
					t.Errorf("unsupported content-type")
			}
		})
	}
}