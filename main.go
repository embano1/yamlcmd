package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

func main() {

	// --- variable definitions
	var (
		input     inputs
		filename  string
		out       []byte
		parallel  int
		wg        sync.WaitGroup
		sudo      bool
		args, exe string
	)

	// --- options for our program
	flag.StringVar(&filename, "f", "config.yaml", "Input file (yaml)")
	flag.IntVar(&parallel, "p", 3, "Number of parallel cmd executions (max: 5)")
	flag.Parse()

	// --- protect ourselves
	if parallel > 5 || parallel < 1 {
		fmt.Printf("%d not allowed for flag \"-p\", using default (3)\n\n", parallel)
		parallel = 5
	}
	token := newSemaphore(parallel)

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

		if input.Ins[k].Type == "command" {

			switch input.Ins[k].Spec.Sudo {
			case true:
				exe = strings.Join([]string{"sudo", input.Ins[k].Spec.ExecCmd}, " ")
				// force OS to cache sudo password so we´re not prompted during command execution
				if !sudo {
					fmt.Printf("-> %s requires sudo \n\n", input.Ins[k].Spec.ExecCmd)
					_, err = exec.Command("/bin/sh", "-c", "sudo -v").CombinedOutput()
					if err != nil {
						fmt.Fprintln(os.Stderr, string(out))
						os.Exit(1)

					} else {
						sudo = true
					}
				}
			default:
				exe = input.Ins[k].Spec.ExecCmd
			}

			args = strings.Join(input.Ins[k].Spec.Args, " ")

			wg.Add(1)
			go func(exe, args string, wg *sync.WaitGroup) {

				// --- acquire token
				<-token.slots

				defer wg.Done()
				out, err = exec.Command("/bin/sh", "-c", exe+" "+args).CombinedOutput()
				if err != nil {
					fmt.Fprintln(os.Stderr, string(out))
				} else {
					//fmt.Printf(")
					fmt.Printf("-> Output of %s %q:\n%s\n", exe, args, string(out))
				}

				// --- release token
				token.slots <- struct{}{}
			}(exe, args, &wg)
		}
	}

	wg.Wait()
	fmt.Println("-> We´re done.")
}

func newSemaphore(permit int) *semaphore {
	sem := &semaphore{slots: make(chan struct{}, permit)}

	for i := 0; i < permit; i++ {
		sem.slots <- struct{}{}
	}
	return sem
}
