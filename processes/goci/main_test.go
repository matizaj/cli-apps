package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		out    string
		expErr error
		setupGit bool
	}{
		{name: "success", proj: "./testdata/tool/", out: "Go Build: SUCCESS\nGo Test: SUCCESS\nGoFmt: SUCCESS\nGit Push: SUCCESS\n", expErr: nil, setupGit: true},
		{name: "fail", proj: "./testdata/toolErr/", out: "failed", expErr: &stepErr{step: "go build"}, setupGit: false},
		{name: "fail format", proj: "./testdata/toolFmtErr/", out: "", expErr: &stepErr{step: "go fmt"}, setupGit:false},
	}

	_, err := exec.LookPath("git")
	if err != nil {
		t.Skip("Git not installed", err)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

			if tc.setupGit {
				cleanup := setupGit(t, tc.proj)
				defer cleanup()
			}
		})
	}
}



func setupGit(t *testing.T, proj string) func() {
	t.Helper()
	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	tempDir, err := os.CreateTemp("", "gocitest")
	if err != nil {
		t.Fatal(err)
	}
	defer tempDir.Close()
	tempDirName := tempDir.Name()

	projPath, err := filepath.Abs(proj)
	if err != nil {
		t.Fatal(err)
	}

	remoteURI := fmt.Sprintf("file://%s", tempDirName)

	var gitCmdList = []struct{
		args []string
		dir string
		env []string
	}{
		{[]string{"init", "--bare"}, tempDirName, nil},
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
		os.RemoveAll(tempDirName)
		os.RemoveAll(filepath.Join(projPath, ".git"))
	}
}