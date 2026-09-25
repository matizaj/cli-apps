package pomodoro_test

import (
	"matizaj/cli-apps/interactveTools/pomo/pomodoro"
	"matizaj/cli-apps/interactveTools/pomo/pomodoro/repository"
	"testing"
	"time"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	return repository.NewInMemeroRepo(), func(){}
}

func TestNewConfig(t *testing.T) {
	testCases:=[]struct{
		name string
		input [3]time.Duration
		expect pomodoro.IntervalConfig
	}{
		{   name: "Default",
			input: [3]time.Duration{
				20*time.Minute,
			}, 
			expect: pomodoro.IntervalConfig{
				PomodoroDuration: 20*time.Minute, 
				ShortBreakDuration: 5*time.Minute, 
				LongBreakDuration: 15*time.Minute,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var repo pomodoro.Repository
			cfg:= pomodoro.NewConfig(repo, tc.input[0], tc.input[1], tc.input[2])

			if cfg.PomodoroDuration != tc.expect.PomodoroDuration {
				t.Errorf("expected pomodoro duration %q got%q", tc.expect.PomodoroDuration, cfg.PomodoroDuration)
			}
		})
	}
}