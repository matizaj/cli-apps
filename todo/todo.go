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

func (l List) String() string {
	formatted:= ""

	for k, t := range l {
		prefix:= "  "
		if t.Done {
			prefix= "X "
		}

		formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
	}

	return formatted
}

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

func (l *List) Verbose(filename string) (string, error) {
	var result string
	if err :=l.Get(filename); err != nil {
		return "",err
	}

	for _, item := range *l {
		result += fmt.Sprintf("Task: %s, Created at: %v, Done: %v\n", item.Task, item.CreatedAt.Format("02-01-2006"), item.Done)
	}
	return result, nil
} 