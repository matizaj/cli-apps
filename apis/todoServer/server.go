package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func newMux(todoFile string) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", rootHandler)
	return m
}

func replyJsonContent(w http.ResponseWriter, r *http.Request, status int, resp *todoResponse) {
	body, err := json.Marshal(resp)
	if err!=nil {
		replyError(w,r,http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func replyError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	log.Printf("%s %s: Error: %d %s", r.URL, r.Method, status, msg)
	http.Error(w, http.StatusText(status), status)
}