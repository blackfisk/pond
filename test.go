package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	runCatFromStdinWorks(populateStdin("testjk:\n"))
	runCatFromStdinWorks(populateStdin("bbb\n"))
}

func populateStdin(str string) func(io.WriteCloser) {
	return func(stdin io.WriteCloser) {
		defer stdin.Close()
		io.Copy(stdin, bytes.NewBufferString(str))
	}
}

func runCatFromStdinWorks(populate_stdin_func func(io.WriteCloser)) {
	cmd := exec.Command("cat")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Panic(err)
	}
	populate_stdin_func(stdin)
	io.Copy(os.Stdout, stdout)
	err = cmd.Wait()
	if err != nil {
		log.Panic(err)
	}
}
