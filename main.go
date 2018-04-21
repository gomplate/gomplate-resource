package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/hairyhenderson/gomplate"
)

// main - check is a no-op for the gomplate resource - the input is always
// considered out-of-date
func main() {
	switch basename := path.Base(os.Args[0]); basename {
	case "check":
		fmt.Printf("[{\"date\": \"%s\"}]\n", time.Now().Format(time.RFC3339))
	case "in":
		in, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("couldn't read stdin: %#v", err)
		}
		var p payload
		err = json.Unmarshal(in, &p)
		if err != nil {
			log.Fatalf("couldn't unmarshal payload: %#v (payload was '%s')", err, string(in))
		}

		var destination string
		if len(os.Args) > 1 {
			destination = os.Args[1]
		}

		p.Params.OutputDir = path.Join(destination, p.Params.OutputDir)
		if len(p.Params.OutputFiles) > 0 {
			for i, f := range p.Params.OutputFiles {
				p.Params.OutputFiles[i] = path.Join(destination, f)
			}
		}

		buf := &bytes.Buffer{}
		gomplate.Stdout = &nopWCloser{buf}
		oStdout := os.Stdout
		oStderr := os.Stderr
		defer func() {
			os.Stdout = oStdout
			os.Stderr = oStderr
		}()
		ro, wo, _ := os.Pipe()
		re, we, _ := os.Pipe()
		os.Stdout = wo
		os.Stdout = we

		stdout := make(chan string)
		stderr := make(chan string)
		go func() {
			b := bytes.Buffer{}
			e := bytes.Buffer{}
			io.Copy(&b, ro)
			io.Copy(&e, re)
			stdout <- b.String()
			stderr <- e.String()
		}()
		wd, _ := os.Getwd()
		err = gomplate.RunTemplates(p.Params)
		wo.Close()
		we.Close()
		os.Stdout = oStdout
		os.Stderr = oStderr
		r := result{
			Version: p.Version,
			Metadata: []metadata{
				{"workDir", wd},
				{"destination", destination},
				{"templateOut", buf.String()},
				{"stdout", <-stdout},
				{"stderr", <-stderr},
				{"success", strconv.FormatBool(err == nil)},
				{"errors", strconv.Itoa(gomplate.Metrics.Errors)},
				{"gatherDuration", gomplate.Metrics.GatherDuration.String()},
				{"totalRenderDuration", gomplate.Metrics.TotalRenderDuration.String()},
			},
		}
		if err != nil {
			r.Metadata = append(r.Metadata, metadata{"error", err.Error()})
		}
		json.NewEncoder(os.Stdout).Encode(r)
		if err != nil {
			log.Fatalf("couldn't run gomplate: %s", err)
		}
	default:
		log.Fatalf("%s is an invalid binary name", basename)
	}
}

// input
type payload struct {
	Version version `json:"version"`
	Source  struct{}
	Params  *gomplate.Config
}

// output
type result struct {
	Version  version    `json:"version"`
	Metadata []metadata `json:"metadata,omitempty"`
}

type version map[string]string

type metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// like ioutil.NopCloser(), except for io.WriteClosers...
type nopWCloser struct {
	io.Writer
}

func (n *nopWCloser) Close() error {
	return nil
}
