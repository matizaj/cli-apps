package main

import (
	"errors"
	"fmt"
	"io"
	"matizaj/cli-apps/todo"
	"net/http"
	"strconv"
	"sync"

)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalidData = errors.New("invalid input")
)
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		replyError(w,r,http.StatusNotFound, "")
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

func validateId(path string, list *todo.List) (int, error) {
	id, err := strconv.Atoi(path)
	if err!= nil {
		return 0, fmt.Errorf("%w: invalid Id: %s", ErrInvalidData, err)
	}

	if id <1 || id > len(*list){
		return 0, fmt.Errorf("%w: invalid Id: %s", ErrInvalidData, err)
	}
	return id, nil
}

func addHandler(w http.ResponseWriter, r *http.Request, list *todo.List, todoFile string) {
	
	todo, err := io.ReadAll(r.Body)
	if err!= nil {
		replyError(w,r,http.StatusInternalServerError, err.Error())
		return
	}
	list.Add(string(todo))
	
	if err := list.Save(todoFile); err!=nil {
		replyError(w,r,http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w,r,http.StatusNoContent, "")
}
func patchHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int, todoFile string) {
	q := r.URL.Query()
	if _, ok := q["complete"]; !ok {
		message := "Missing query param 'complete'"
		replyError(w,r,http.StatusBadRequest, message)
		return
	}

	list.Complete(id)
	if err := list.Save(todoFile); err!=nil {
		replyError(w,r,http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w,r,http.StatusNoContent, "")
}

func deleteHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int, todoFile string) {
	list.Delete(id)
	if err := list.Save(todoFile); err!=nil {
		replyError(w,r,http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w,r,http.StatusNoContent, "")
}

func getOneHandler(w http.ResponseWriter, r *http.Request, list *todo.List, id int) {
	resp:=&todoResponse{
		Results: (*list)[id-1:id],
	}
	replyJsonContent(w,r,http.StatusOK, resp)
}

func getAllHandler(w http.ResponseWriter, r *http.Request, list *todo.List) {
	resp:=&todoResponse{
		Results: *list,
	}
	replyJsonContent(w,r,http.StatusOK, resp)
}