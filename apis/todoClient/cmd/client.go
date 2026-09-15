package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrConnection = errors.New("connection error")
	ErrNotForund = errors.New("not found")
	ErrInvalidResponse = errors.New("invalid response")
	ErrInvalid = errors.New("invalid data")
	ErrNotNumber = errors.New("not a number")
)

type item struct {
	Task string
	Done bool
	CreatedAt time.Time
	CompletedAt time.Time
}

type response struct {
	Results []item	`json:results`
	Date int	`json:date`
	Totalresults int	`json:total_results`
}

func newClient() *http.Client {
	return &http.Client{
		Timeout: 10*time.Second,
	}
}

func getItems(url string) ([]item, error) {
	r, err := newClient().Get(url)
	if err !=nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}

	defer r.Body.Close()

	if r.StatusCode!= http.StatusOK {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to process request %s", string(msg))
		}
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err=ErrNotForund
		}
		return nil, fmt.Errorf("%w: %s", err, msg)
		
	}

	var resp response
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}

	if resp.Totalresults ==0 {
		return nil, fmt.Errorf("%w", ErrNotForund)
	}

	return resp.Results, nil
}

func getAll(apiroot string) ([]item, error) {
	u:=fmt.Sprintf("%s/todo", apiroot)
	return getItems(u)
}

func getOne(url string, id int) (item, error) {
	u:=fmt.Sprintf("%s/todo/%d?complete", url, id)
	r, err := newClient().Get(u)
	if err != nil {
		return item{}, err
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return item{}, fmt.Errorf("%w: %s", ErrNotForund, err)
	}

	var resp response

	if err :=json.NewDecoder(r.Body).Decode(&resp);err!= nil {
		return item{}, err
	}

	return resp.Results[0], nil
}

func addItem(hosturl string, task string) error {
	u:=fmt.Sprintf("%s/todo", hosturl)
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(task))
	if err!= nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err!= nil {
		return err
	}
	defer resp.Body.Close()
	
	fmt.Println(resp)
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("%w: %s", ErrInvalidResponse, err)
	}
	return nil
}

func completeItem(hosturl string, id int) error {
	url, err:=newClient().Get(hosturl)
	if err!= nil {
		return err
	}
	u:=fmt.Sprintf("%s/todo/%d", url, id)
	req, err := http.NewRequest(http.MethodPatch, u, nil)
	if err!= nil {
		return err
	}

	defer req.Body.Close()

	r, err := http.DefaultClient.Do(req)
	if err!= nil {
		return err
	}

	if r.StatusCode != http.StatusNoContent {
		return fmt.Errorf("not completed %w", err)
	}
	return nil
}