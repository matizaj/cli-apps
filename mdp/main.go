package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
		<html>
		<head>
			<meta http-equiv="content-type" content="text/html;charset=utf-8">
		<title>Markdown Preview Tool</title>
		</head>
		<body>`

	footer = `
	</body>
	</html>`
)


func main() {
	fmt.Println("MardDown Preview")

	filename:=flag.String("file", "", "Markdown file to preview")
	flag.Parse()

	if *filename == "" {
		fmt.Fprintln(os.Stderr, "no file provided")
		os.Exit(1)
	}

	if err:=run(*filename); err!= nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}
}

func run(filename string) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}

	htmlData := parseContent(bytes)

	outName:=fmt.Sprintf("%s.html", filepath.Base(filename))
	fmt.Print(outName)

	return saveHtml(outName, htmlData)
}

func parseContent(b []byte) []byte {
	var buffer bytes.Buffer
	output:= blackfriday.Run(b)

	body:= bluemonday.UGCPolicy().SanitizeBytes(output)

	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)
	return buffer.Bytes()
}

func saveHtml(outName string, html []byte) error {
	return os.WriteFile(outName, html, 0644)
}