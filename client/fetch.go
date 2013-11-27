package client

import (
	"code.google.com/p/gopass"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
)

type PondClient struct {
	Home           string
	url            string
	agentAvailable bool
	passphrase     string
}

func NewPondClient(url string) *PondClient {
	pc := new(PondClient)
	pc.url = url

	return pc
}

func (c *PondClient) getJSON() []interface{} {
	response, _ := http.Get(c.url)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		os.Exit(1)
	}
	var data []interface{}

	if err := json.Unmarshal(contents, &data); err != nil {
		panic(err)
	}

	return data
}

func (c *PondClient) createDir() {
	usr, _ := user.Current()
	home := fmt.Sprintf("%s/%s", usr.HomeDir, ".pond/messages")
	os.MkdirAll(home, 0777)

	c.Home = home
}

func (r *PondClient) messageHash(msg string) string {
	h := sha1.New()
	h.Write([]byte(msg))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return sha
}

func (c *PondClient) decryptMessages(data []interface{}) {
	if !c.agentAvailable {
		c.passphrase, _ = gopass.GetPass("Enter the passphrase: ")
	}

	for _, message := range data {
		msg := message.(string)
		sha := c.messageHash(msg)

		gpg_filename := fmt.Sprintf("%s/%s.gpg", c.Home, sha)
		filename := fmt.Sprintf("%s/%s", c.Home, sha)

		err := ioutil.WriteFile(gpg_filename, []byte(msg), 0666)
		if err != nil {
			panic(err)
		}

		var args []string

		if !c.agentAvailable {
			args = []string{
				"--batch",
				"--passphrase", c.passphrase,
				"-o", filename, "--decrypt", gpg_filename}
		} else {
			args = []string{
				"--batch", "--use-agent",
				"-o", filename, "--decrypt", gpg_filename}
		}

		cmd := exec.Command("gpg", args...)
		cmd.Start()
		cmd.Wait()
	}
}

func (c *PondClient) readMessages() {
	new_messages, _ := filepath.Glob(c.Home + "/*")
	for _, incoming := range new_messages {
		ext := filepath.Ext(incoming)
		if ext == "" {
			content, _ := ioutil.ReadFile(incoming)
			fmt.Println("-----------------------------------------")
			fmt.Println(string(content))
			fmt.Println("-----------------------------------------")
		}
	}
}

func (c *PondClient) agentIsRunning() {
	out, err := exec.Command("ps", "uax").Output()
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile("gpg-agent")
	running := re.FindString(string(out)) != ""
	c.agentAvailable = running
}

func (c *PondClient) Fetch() {
	data := c.getJSON()

	c.createDir()
	c.agentIsRunning()
	c.decryptMessages(data)
	c.readMessages()
}
