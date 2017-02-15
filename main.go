package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {

	// --- variable definitions
	var input inputs
	var filename string
	var out []byte

	// --- options for our program
	flag.StringVar(&filename, "f", "config.yaml", "Input file (yaml)")
	flag.Parse()

	// --- yaml functions
	ymlsrc, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(ymlsrc, &input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// --- command execution
	for k := range input.Ins {

		var args, exe string

		if input.Ins[k].Type == "command" {

			switch input.Ins[k].Spec.Sudo {
			case true:
				exe = strings.Join([]string{"sudo", input.Ins[k].Spec.ExecCmd}, " ")
			default:
				exe = input.Ins[k].Spec.ExecCmd
			}

			args = strings.Join(input.Ins[k].Spec.Args, " ")
			fmt.Printf("-> Calling %s with argument(s) %q:\n", exe, args)

			out, err = exec.Command("/bin/sh", "-c", exe+" "+args).CombinedOutput()
			if err != nil {
				fmt.Fprintln(os.Stderr, string(out))
				continue
			}

			fmt.Println(string(out))
		}
	}
}
