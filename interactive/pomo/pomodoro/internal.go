package pomodoro

//category const
const (
	CategoryPomodoro="Pomodoro"
	CategoryShortBreak="ShortBreak"
	CategoryLongBreak="LongBreak"
)

//state const
const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)
