package repository

import (
	"fmt"
	"matizaj/cli-apps/interactveTools/pomo/pomodoro"
	"sync"
)

type inMemoryRepo struct {
	sync.RWMutex
	intervals []pomodoro.Interval
}

func NewInMemeroRepo() *inMemoryRepo {
	return &inMemoryRepo{
		intervals: []pomodoro.Interval{},
	}
}

func (r *inMemoryRepo) Create(i pomodoro.Interval)(int64, error) {
	r.Lock()
	defer r.Unlock()

	i.Id = int64(len(r.intervals))+1
	r.intervals = append(r.intervals, i)
	return i.Id, nil
}

func (r *inMemoryRepo)Update(i pomodoro.Interval)error {
	r.Lock()
	defer r.Unlock()

	if i.Id == 0 {
		return fmt.Errorf("%w: %d", pomodoro.ErrInvalidId, i.Id)
	}
	r.intervals[i.Id-1] = i
	return nil
}