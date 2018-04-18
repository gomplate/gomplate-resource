package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/hairyhenderson/gomplate"
)

// main - check is a no-op for the gomplate resource - the input is always
// considered out-of-date
func main() {
	switch basename := path.Base(os.Args[0]); basename {
	case "check":
		fmt.Printf("[{\"noop\":true,\"date\": \"%s\"}]\n", time.Now().Format(time.RFC3339))
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
		err = gomplate.RunTemplates(p.Params)
		if err != nil {
			log.Fatalf("couldn't run gomplate: %#v", err)
		}
	default:
		log.Fatalf("%s is an invalid binary name", basename)
	}
}

// input payload
type payload struct {
	Source struct{}
	Params *gomplate.Config
}
