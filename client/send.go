package client

import (
	"bytes"
	"io"
	"log"
	"errors"
	"net/http"
	"os/exec"
)

func (c *PondClient) populateStdin(str string) func(io.WriteCloser) {
	return func(stdin io.WriteCloser) {
		defer stdin.Close()
		io.Copy(stdin, bytes.NewBufferString(str))
	}
}

func (c *PondClient) Send(email, message string) bool {
	gpg, err := c.runCmdFromStdin(c.populateStdin(message), email)
        if err != nil {
                return false
        }
	return c.sendToPond(gpg)
}

func (c *PondClient) sendToPond(message string) bool {
	buf := &bytes.Buffer{}
	buf.Write([]byte(message))
	response, _ := http.Post(c.url, "text/plain", buf)
	defer response.Body.Close()

	return response.StatusCode == 200
}

func (c *PondClient) runCmdFromStdin(populate_stdin_func func(io.WriteCloser), email string) (string, error) {
	args := []string{
		"--encrypt", "--armor",
		"-R", email,
	}

	buf := &bytes.Buffer{}
	cmd := exec.Command("gpg", args...)
	cmd.Stdout = buf

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Panic(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Panic(err)
	}
	populate_stdin_func(stdin)

	err = cmd.Wait()
	if err != nil {
                return "", errors.New("Invalid email")
	}
	return buf.String(), nil
}
