package pomodoro

import (
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