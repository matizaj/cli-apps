package repository

import (
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