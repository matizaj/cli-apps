package cmd

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

func TestAddAction(t *testing.T) {
	expUrl:="/todo"
	expMethod:= http.MethodPost
	expBody:="{\"task\":\"Task 1\"}\n"
	expContentType:="application/json"
	args:="Task 1"

	url, cleanup:=mockServer(func(w http.ResponseWriter, r * http.Request) {
		if r.URL.Path != expUrl {
			t.Errorf("expected path %q got %q",expUrl, r.URL.Path)
		}
		if r.Method != expMethod {
			t.Errorf("expected method type %q got %q",expMethod, r.Method)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}

		defer r.Body.Close()

		if string(body) != expBody {
			t.Errorf("expected body %q got %q",expBody, string(body))
		}
		ct:= r.Header.Get("Content-Type")
		if ct!= expContentType {
			t.Errorf("expected content-type %q got %q",expContentType, ct)
		}
		w.WriteHeader(testResp["created"].Status)
		
	})

	defer cleanup()

	var out bytes.Buffer
	if err := addAction(&out, url, args); err != nil {
		t.Fatalf("expected no err, got %q", err)
	}
}