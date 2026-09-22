package pomodoro

import (
	"context"
	"errors"
	"time"
)

//category const
const (
	CategoryPomodoro   = "Pomodoro"
	CategoryShortBreak = "ShortBreak"
	CategoryLongBreak  = "LongBreak"
)

//state const
const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

type Interval struct {
	Id              int64
	StartTime       time.Time
	PlannedDuration time.Duration
	ActualDuration  time.Duration
	Category        string
	State           int
}

type Repository interface {
	Create(i Interval)(int64, error)
	ById(id int64)(Interval, error)
	Breaks(id int64)([]Interval, error)
	Last()(Interval, error)
	Update(i Interval)error
}

var (
	ErrNoIntervals = errors.New("no intervals")
	ErrIntervalNotRunning = errors.New("interval not running")
	ErrNoIntervalCompleted = errors.New("interval is completed or cancelled")
	ErrInvalidState = errors.New("invalid state")
	ErrInvalidId = errors.New("invalid id")
)

type IntervalConfig struct {
	repo Repository
	PomodoroDuration time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration time.Duration
}

type Callback func(Interval)


func NewConfig(repo Repository, pomodoro, shortBreak, longBreak time.Duration) *IntervalConfig{
	cfg:= &IntervalConfig{
		repo: repo,
		PomodoroDuration: 25*time.Minute,
		ShortBreakDuration: 5*time.Minute,
		LongBreakDuration: 15*time.Minute,
	}

	if pomodoro>0{
		cfg.PomodoroDuration=pomodoro
	}
	if shortBreak>0{
		cfg.ShortBreakDuration=shortBreak
	}
	if longBreak>0{
		cfg.LongBreakDuration=longBreak
	}

	return cfg
}

func nextCategory(r Repository)(string, error) {
	li, err := r.Last()
	if err != nil && err == ErrNoIntervals{
		return CategoryPomodoro, err
	}
	if err != nil {
		return "", err
	}

	if li.Category == CategoryLongBreak || li.Category == CategoryShortBreak {
		return CategoryPomodoro, nil
	}

	lastBreaks, err := r.Breaks(3)
	if err != nil {
		return "", err
	}

	if len(lastBreaks) <3 {
		return CategoryShortBreak, nil
	}

	for _, i := range lastBreaks {
		if i.Category == CategoryLongBreak {
			return CategoryShortBreak, nil
		}
	}

	return CategoryLongBreak, nil
}

func tick(ctx context.Context, id int64, config *IntervalConfig, start, periodic, end Callback) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	i, err := config.repo.ById(id)
	if err != nil {
		return err
	}

	expire := time.After(i.PlannedDuration - i.ActualDuration)
	start(i)

	for {
		select {
		case <-ticker.C:
			i, err := config.repo.ById(id)
			if err != nil {
				return err
			}

			if i.State == StatePaused {
				return nil
			}
			i.ActualDuration += time.Second
			if err:= config.repo.Update(i); err != nil {
				return err
			}

			periodic(i)
		case <-expire:
			i, err := config.repo.ById(id)
			if err != nil {
				return err
			}
			i.State= StateDone
			end(i)
			return config.repo.Update(i)
		case <-ctx.Done():
			i, err := config.repo.ById(id)
			if err != nil {
				return err
			}
			i.State = StateCancelled
			return config.repo.Update(i)
		}
	}
}

func newInterval(cfg *IntervalConfig) (Interval, error) {
	i := Interval{}

	cat, err := nextCategory(cfg.repo)
	if err != nil {
		return i, err
	}
	i.Category = cat
	switch cat {
	case CategoryPomodoro:
		i.PlannedDuration = cfg.PomodoroDuration
	case CategoryShortBreak:
		i.PlannedDuration = cfg.ShortBreakDuration
	case CategoryLongBreak:
		i.PlannedDuration = cfg.LongBreakDuration
	}

	if i.Id, err = cfg.repo.Create(i); err != nil {
		return i, err
	}

	return i, nil 
}