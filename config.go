package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func editConfig() error {
	f, err := getConfigPath("hosts")
	if err != nil {
		return err
	}
	fname := f.Name()
	writeExample(f)
	defer f.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, fname)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return errors.New("Editor " + editor + " not found. Please update your $EDITOR environment variable")
	}
	err = cmd.Wait()

	return err
}

func readHostConfig() ([]string, error) {
	var hosts []string
	f, err := getConfigPath("hosts")
	if err != nil {
		return hosts, err
	}
	defer f.Close()

	hosts, err = parseConfig(f)
	if len(hosts) == 0 {
		f.Close()
		editConfig()
		f, _ := getConfigPath("hosts")
		hosts, err = parseConfig(f)
		f.Close()
		if len(hosts) == 0 {
			return hosts, errors.New("Servers config can't be empty")
		}
	}
	return hosts, err
}

var example = `# Put host config to be line separated. example:
# root@192.168.0.123
# root@192.168.0.124
`

func writeExample(f *os.File) {
	b, _ := ioutil.ReadAll(f)
	if len(b) == 0 {
		w := bufio.NewWriter(f)
		_, err := w.WriteString(example)
		err = w.Flush()
		if err != nil {
			log.Println(err)
		}
	}
}

func parseConfig(f *os.File) ([]string, error) {
	var hosts []string
	byt, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(byt), "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)

		if strings.HasPrefix(l, "#") || l == "" {
			continue
		}
		hosts = append(hosts, l)
	}
	return hosts, nil
}

func getConfigPath(file string) (*os.File, error) {
	path := os.Getenv("HOME")
	path += "/.config/rexec/"
	err := os.MkdirAll(path, 0776)
	if err != nil {
		return nil, err
	}

	fpath := path + file
	f, err := os.OpenFile(fpath, os.O_RDWR, os.ModeDevice)
	if os.IsNotExist(err) {
		f, err = os.Create(fpath)
	}

	return f, err
}
