package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{name: "NoFilter", root: "testdata",
			cfg:      config{ext: "", size: 0, list: true},
			expected: filepath.Join("testdata", "dir.log") + "\n" + filepath.Join("testdata", "dir2", "script.sh") + "\n"},
		{name: "FilterExtensionMatch", root: "testdata",
			cfg:      config{ext: ".log", size: 0, list: true},
			expected: filepath.Join("testdata", "dir.log") + "\n"},
		{name: "FilterExtensionSizeMatch", root: "testdata",
			cfg:      config{ext: ".log", size: 10, list: true},
			expected: filepath.Join("testdata", "dir.log") + "\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata",
			cfg:      config{ext: ".log", size: 20, list: true},
			expected: ""},
		{name: "FilterExtensionNoMatch", root: "testdata",
			cfg:      config{ext: ".gz", size: 0, list: true},
			expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()
			if res != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, res)
			}
		})
	}
}

func TestRunDelExtensions(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nNoDelete   int
		nDelete     int
		expected    string
	}{
		{name: "DeleteExtensionNoMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz", nDelete: 0, nNoDelete: 10,
			expected: ""},
		{name: "DeleteExtensionMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: "", nDelete: 10, nNoDelete: 0,
			expected: ""},
		{name: "DeleteExtensionMixed",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz", nDelete: 5, nNoDelete: 5,
			expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)
			tc.cfg.wLog = &logBuffer

			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})

			defer cleanup()
			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}
			res := buffer.String()

			if tc.expected != res {
				t.Errorf("expected %q, got %q", tc.expected, res)
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("expected %d files left, got %d", tc.nNoDelete, len(filesLeft))
			}
			expLogLines := tc.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))

			if len(lines) != expLogLines {
				t.Errorf("Expected %d log lines, got %d instead\n",
					expLogLines, len(lines))
			}
		})
	}
}

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()
	tempDir, _ := os.MkdirTemp("", "walk-test-*")
	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDir, fname)
			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return tempDir, func() { os.RemoveAll(tempDir) }
}

func TestRunArchive(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          config
		extNoArchive string
		nNoArchive   int
		nArchive     int
	}{
		{name: "ArchiveExtensionNoMatch", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 0, nNoArchive: 10},
		{name: "ArchiveExtensionMatch", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 10, nNoArchive: 0},
		{name: "ArchiveExtensionNoMixed", cfg: config{ext: ".log"}, extNoArchive: ".gz", nArchive: 5, nNoArchive: 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:      tc.nArchive,
				tc.extNoArchive: tc.nNoArchive,
			})

			defer cleanup()

			archiveDir, cleanupArchive := createTempDir(t, nil)
			defer cleanupArchive()

			tc.cfg.archive = archiveDir

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			if got := buffer.String(); got != "" {
				t.Errorf("expected no standard output while archiving, got %q", got)
			}
			filesArchived, err := os.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(filesArchived) != tc.nArchive {
				t.Errorf("expected %d, got %d\n", tc.nArchive, len(filesArchived))
			}

		})
	}
}
