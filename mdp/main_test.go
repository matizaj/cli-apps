package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

)

const (
	inputFile  = "./testdata/test.md"
	resultFile = "test.md.html"
	goldenFile = "./testdata/test.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	result := parseContent(input)
	// os.WriteFile(goldenFile, result, 0644) - make sure golden is a golden :)
	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden %s\n", expected)
		t.Logf("result %s\n", result)
		t.Error("result content does not match1")
	}
}

func TestRun(t *testing.T) {
	var mockStdout bytes.Buffer
	if err := run(inputFile, &mockStdout); err != nil {
		t.Fatal(err)
	}
	resultFile := strings.TrimSpace(mockStdout.String())
	result, err:= os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err:= os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden %s\n", expected)
		t.Logf("result %s\n", result)
		t.Error("result content does not match1")
	}

	os.Remove(resultFile)

}