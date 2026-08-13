package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>{{ .Title }}</title>
  </head>
  <body>
{{ .Body }}
  </body>
</html>
`
)

type content struct {
	Title string
	Body template.HTML
}


func main() {
	fmt.Println("MardDown Preview")

	filename:=flag.String("file", "", "Markdown file to preview")
	tFname:=flag.String("t", "", "Alternate template name")
	skipPreview := flag.Bool("s", false, "skip auto-preview")
	flag.Parse()

	if *filename == "" {
		fmt.Fprintln(os.Stderr, "no file provided")
		os.Exit(1)
	}

	if err:=run(*filename, *tFname, os.Stdout, *skipPreview); err!= nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}
}

func run(filename string, tFname string, w io.Writer, skipPreview bool) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}

	htmlData, err := parseContent(bytes, tFname)
	// todo: err handlid 

	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}
	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()
	
	fmt.Fprintln(w, outName)

	if err:= saveHtml(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}
	
	defer os.Remove(outName)

	return preview(outName)
}

func parseContent(b []byte, tFname string) ([]byte, error) {
	var buffer bytes.Buffer
	output:= blackfriday.Run(b)

	body:= bluemonday.UGCPolicy().SanitizeBytes(output)
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
		return nil, err
		}
	}

	c := content {
		Title: "Markdown Preview Tool",
		Body: template.HTML(body),
	}
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}
	// buffer.WriteString(header)
	// buffer.Write(body)
	// buffer.WriteString(footer)
	return buffer.Bytes(), nil
}

func saveHtml(outName string, html []byte) error {
	return os.WriteFile(outName, html, 0644)
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName="cmd.exe"
		cParams=[]string{"/C","start"}
	case "darwin":
		cName= "open"
	default:
		return fmt.Errorf("OS not supported")
	}
	cParams = append(cParams, fname)

	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	err = exec.Command(cPath, cParams...).Run()
	time.Sleep(2*time.Second)
	return err
}