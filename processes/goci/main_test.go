package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)
func TestRunKill(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		sig    syscall.Signal
		expErr error
	}{
		{"SIGINT", "./testdata/tool", syscall.SIGINT, ErrSignal},
		{"SIGTERM", "./testdata/tool", syscall.SIGTERM, ErrSignal},
		{"SIGQUIT", "./testdata/tool", syscall.SIGQUIT, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			command=mockCmdTimeout
			errChan:=make(chan error)
			ignSigChan:=make(chan os.Signal,1)
			expSigChan:=make(chan os.Signal,1)

			signal.Notify(ignSigChan, syscall.SIGQUIT)
			defer signal.Stop(ignSigChan)

			signal.Notify(expSigChan, tc.sig)
			defer signal.Stop(expSigChan)

			go func(){
				errChan<- run(tc.proj, io.Discard)
			}()

			go func() {
				time.Sleep(2*time.Second)
				p, err := os.FindProcess(os.Getpid())
				if err == nil {
					err = p.Signal(tc.sig)
				}
			}()
			select {
			case err := <- errChan:
				if err == nil {
					t.Errorf("expected err but got nil")
					return
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected %q, but got %q", tc.expErr, err)
				}

				// select signal
				select {
					case rec := <-expSigChan:
					if rec != tc.sig {
						t.Errorf("Expected signal %q, got %q", tc.sig, rec)
					}
					default:
					t.Errorf("Signal not received")
					}
				case <-ignSigChan:
			}
		})
	}
}
func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		out    string
		expErr error
		setupGit bool
		mockCmd func(ctx context.Context, exec string, args ...string) *exec.Cmd
	}{
		{name: "success", proj: "./testdata/tool/", out: "Go Build: SUCCESS\nGo Test: SUCCESS\nGoFmt: SUCCESS\nGit Push: SUCCESS\n", expErr: nil, setupGit: true, mockCmd: nil},
		{name: "success mock", proj: "./testdata/tool/", out: "Go Build: SUCCESS\nGo Test: SUCCESS\nGoFmt: SUCCESS\nGit Push: SUCCESS\n", expErr: nil, setupGit: true, mockCmd: mockCmdContext},
		{name: "fail", proj: "./testdata/toolErr/", out: "failed", expErr: &stepErr{step: "go build"}, setupGit: false, mockCmd: nil},
		{name: "fail format", proj: "./testdata/toolFmtErr/", out: "", expErr: &stepErr{step: "go fmt"}, setupGit:false, mockCmd: nil},
		{name: "fail timeout", proj: "./testdata/tool/", out: "", expErr: context.DeadlineExceeded, setupGit:false, mockCmd: mockCmdTimeout},
	}



	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupGit {
				_, err := exec.LookPath("git")
				if err != nil {
					t.Skip("Git not installed", err)
				}
				cleanup:=setupGit(t, tc.proj)
				defer cleanup()
			}
			
			if tc.mockCmd!= nil {
				command = tc.mockCmd
			}
		
			var buffer bytes.Buffer
			err := run(tc.proj, &buffer)

			fmt.Printf("tc: %s, output: %v\n", tc.name, buffer.String())

			if tc.expErr != nil {
				if err == nil {
					t.Errorf("expected err: %q but got 'nil'", tc.expErr)
				}

				if !errors.Is(err, tc.expErr) {
					t.Errorf("expected err: %q, but got %q", tc.expErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected err: %q", err)
			}

			if buffer.String() != tc.out {
				t.Errorf("expected output %s, but got %s", tc.out, buffer.String())
			}
		})
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS")!="1"{
		return
	}
	if os.Getenv("GO_HELPER_TIMEOUT")=="1"{
		time.Sleep(5*time.Second  )
	}
	if os.Args[2] == "git" {
		fmt.Fprintln(os.Stdout, "everything up-to-date")
		os.Exit(0)
	}
	os.Exit(1)
}



func setupGit(t *testing.T, proj string) func() {
	t.Helper()
	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	tempDir, err := os.MkdirTemp("", "gocitest")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("temp dir name: %s\n", tempDir)

	projPath, err := filepath.Abs(proj)
	if err != nil {
		t.Fatal(err)
	}

	remoteURI := fmt.Sprintf("file://%s", tempDir)

	var gitCmdList = []struct{
		args []string
		dir string
		env []string
	}{
		{[]string{"init", "--bare"}, tempDir, nil},
		{[]string{"init"}, projPath, nil},
		{[]string{"remote", "add", "origin",remoteURI},projPath , nil},
		{[]string{"add", "."},projPath , nil},
		{[]string{"commit", "-m", "test"},projPath , []string{
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@example.com",
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@example.com",
		}},
	}

	for _, g := range gitCmdList {
		gitCmd:= exec.Command(gitExec, g.args...)
		gitCmd.Dir=g.dir

		if g.env != nil {
			gitCmd.Env = append(os.Environ(), g.env...)
		}

		if err := gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}

	return func() {
		os.RemoveAll(tempDir)
		os.RemoveAll(filepath.Join(projPath, ".git"))
	}
}

func mockCmdContext(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess"}
	cs = append(cs, exe)
	cs = append(cs, args...)
	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func mockCmdTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess"}
	cs = append(cs, exe)
	cs = append(cs, args...)
	cmd := mockCmdContext(ctx, exe, args...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}