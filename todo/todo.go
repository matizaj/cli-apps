package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task      string
	Done      bool
	CreatedAt time.Time
	CompletedAt time.Time
}

type List []item

func (l *List) Add(task string) {
	item := item {
		Task: task,
		CreatedAt: time.Now(),
	}

	*l = append(*l, item)
}

func (l *List) Complete(i int) error {
	list := *l
	if i<=0 || i> len(list) {
		return fmt.Errorf("item %d does not exist\n", i)
	}

	list[i-1].Done = true
	list[i-1].CompletedAt = time.Now()

	return nil
}


func (l *List) Delete(i int) error {
	list := *l
	if i<=0 || i> len(list) {
		return fmt.Errorf("item %d does not exist\n", i)
	}
	*l = append(list[:i-1], list[i:]...)

	return nil
}

func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)
}