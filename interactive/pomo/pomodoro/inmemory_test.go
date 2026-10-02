package pomodoro_test

import (
	"context"
	"errors"
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
			expect: pomodoro.IntervalConfig{
				PomodoroDuration: 25*time.Minute, 
				ShortBreakDuration: 5*time.Minute, 
				LongBreakDuration: 15*time.Minute,
			},
		},
		{   name: "SingleInput",
			input: [3]time.Duration{
				20*time.Minute,
			}, 
			expect: pomodoro.IntervalConfig{
				PomodoroDuration: 20*time.Minute, 
				ShortBreakDuration: 5*time.Minute, 
				LongBreakDuration: 15*time.Minute,
			},
		},
		{   name: "MultiInput",
			input: [3]time.Duration{
				20*time.Minute,
				10*time.Minute,
				12*time.Minute,
			}, 
			expect: pomodoro.IntervalConfig{
				PomodoroDuration: 20*time.Minute, 
				ShortBreakDuration: 10*time.Minute, 
				LongBreakDuration: 12*time.Minute,
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

func TestStart(t *testing.T) {
  const duration = 2 * time.Second

  repo, cleanup := getRepo(t)
  defer cleanup()

  config := pomodoro.NewConfig(repo, duration, duration, duration)

  testCases := []struct {
    name        string
    cancel      bool
    expState    int
    expDuration time.Duration
  }{

    {name: "Finish", cancel: false,
      expState: pomodoro.StateDone, expDuration: duration},
    {name: "Cancel", cancel: true,
      expState: pomodoro.StateCancelled, expDuration: duration / 2},
  }

  // Execute tests for Start
  for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
      ctx, cancel := context.WithCancel(context.Background())

      i, err := pomodoro.GetInterval(config)
      if err != nil {
        t.Fatal(err)
      }

      start := func(i pomodoro.Interval) {
        if i.State != pomodoro.StateRunning {
          t.Errorf("Expected state %d, got %d.\n",
            pomodoro.StateRunning, i.State)
        }
        if i.ActualDuration >= i.PlannedDuration {
          t.Errorf("Expected ActualDuration %q, less than Planned %q.\n",
            i.ActualDuration, i.PlannedDuration)
        }
      }

      end := func(i pomodoro.Interval) {
        if i.State != tc.expState {
          t.Errorf("Expected state %d, got %d.\n",
            tc.expState, i.State)
        }
        if tc.cancel {
          t.Errorf("End callback should not be executed")
        }
      }

      periodic := func(i pomodoro.Interval) {
        if i.State != pomodoro.StateRunning {
          t.Errorf("Expected state %d, got %d.\n",
            pomodoro.StateRunning, i.State)
        }
        if tc.cancel {
          cancel()
        }
      }
      
      if err := i.Start(ctx, config, start, periodic, end); err != nil {
        t.Fatal(err)
      }

      i, err = repo.ById(i.Id)
      if err != nil {
        t.Fatal(err)
      }

      if i.State != tc.expState {
        t.Errorf("Expected state %d, got %d.\n",
          tc.expState, i.State)
      }
      if i.ActualDuration != tc.expDuration {
        t.Errorf("Expected ActualDuration %q, got %q.\n",
          tc.expDuration, i.ActualDuration)
      }
      cancel()
    })
  }
}

func TestPause(t *testing.T) {
  const duration = 2 * time.Second

  repo, cleanup := getRepo(t)
  defer cleanup()

  cfg := pomodoro.NewConfig(repo, duration, duration, duration)
  testCases := []struct{
    name string
    start bool
    expState int
    expDuration time.Duration
  }{
    {name: "NotStarted", start: false, expState: pomodoro.StateNotStarted, expDuration: 0},
    {name: "Paused", start: true, expState: pomodoro.StatePaused, expDuration: duration/2},
  }

  expError := pomodoro.ErrIntervalNotRunning

  for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
      ctx, cancel := context.WithCancel(context.Background())

      i, err := pomodoro.GetInterval(cfg)
      if err != nil {
        t.Fatal(err)
      }

      start := func(pomodoro.Interval) {}
      end := func(pomodoro.Interval) {
        t.Errorf("End callback should not be executed")
      }
      periodic := func(i pomodoro.Interval) {
        if err := i.Pause(cfg); err != nil {
          t.Fatal(err)
        }
      }

      if tc.start {
        if err := i.Start(ctx, cfg, start, periodic, end); err != nil {
          t.Fatal(err)
        }
      }

      i, err = pomodoro.GetInterval(cfg)
      if err != nil {
        t.Fatal(err)
      }

      err = i.Pause(cfg)
      if err != nil {
        if ! errors.Is(err, expError) {
          t.Fatalf("Expected error %q, got %q", expError, err)
        }
      }

      if err == nil {
        t.Errorf("Expected error %q, got nil", expError)
      }

      i, err = repo.ById(i.Id)
      if err != nil {
        t.Fatal(err)
      }

      if i.State != tc.expState {
        t.Errorf("Expected state %d, got %d.\n",
          tc.expState, i.State)
      }

      if i.ActualDuration != tc.expDuration {
        t.Errorf("Expected duration %q, got %q.\n",
          tc.expDuration, i.ActualDuration)
      }
      cancel()
    })

  }
  
}
