package main

import (
	"errors"
	"matizaj/cli-apps/todo"
	"net/http"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalidData = errors.New("invalid input")
)
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w,r)
		return
	}
	content := "todo server api"
	replyTextContent(w,r,http.StatusOK, content)
}

func replyTextContent(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(content))
}


func todoRouter(todoFile string, l sync.Locker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := &todo.List{}
		l.Lock()
		defer l.Unlock()

		if err := list.Get(todoFile); err!= nil {
			replyError(w,r,http.StatusInternalServerError, err.Error())
			return
		}
		if r.URL.Path == "" {
			switch r.Method {
			case http.MethodGet:
				getAllHandler(w,r,list)
			case http.MethodPost:
				addHandler(w,r,list, todoFile)
			default:
				message:="Method not supported"
				replyError(w,r,http.StatusMethodNotAllowed, message)
			}
			return 
		}

		id, err := validateId(r.URL.Path, list)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				replyError(w,r,http.StatusNotFound, err.Error())
				return
			}
			replyError(w,r,http.StatusBadRequest, err.Error())
				return
		}

		switch r.Method {
		case http.MethodGet:
			getOneHandler(w,r,list, id)
		case http.MethodDelete:
			deleteHandler(w,r,list, id, todoFile)
		case http.MethodPatch:
			patchHandler(w,r,list,id, todoFile)
		default:
			message:= "Method not supported"
			replyError(w,r,http.StatusMethodNotAllowed, message)
		}		
	}
}