package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
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

func readHostConfig(group string) ([]string, error) {
	var hosts []string
	f, err := getConfigPath("hosts")
	if err != nil {
		return hosts, err
	}
	defer f.Close()

	hosts, err = parseConfig(readFile(f), group)
	if len(hosts) == 0 {
		f.Close()
		editConfig()
		f, _ := getConfigPath("hosts")
		hosts, err = parseConfig(readFile(f), group)
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
#
# or you can group it. example:
# [server1]
# root@192.168.0.1
# root@192.168.0.2
#
# [other-server]
# root@192.168.0.3
# root@192.168.0.4
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

var groupRegex = regexp.MustCompile(`\[([\w\-]+)]`)

func readFile(f *os.File) string {
	byt, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return string(byt)
}

func parseConfig(rawConfig, groupConfig string) ([]string, error) {
	var hosts []string
	var groups = make(map[string][]string)

	var groupMode bool
	var currentGroup string
	lines := strings.Split(rawConfig, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)

		if strings.HasPrefix(l, "#") || l == "" {
			currentGroup = ""
			continue
		}
		if bt := groupRegex.Find([]byte(l)); len(bt) > 0 {
			groupMode = true
			group := groupRegex.ReplaceAllString(l, "$1") // remove [ and ]
			currentGroup = group
			continue
		}

		if currentGroup != "" {
			groups[currentGroup] = append(groups[currentGroup], l)
			continue
		}

		if groupMode {
			return nil, fmt.Errorf("Found config with group mode. but %s doesn't have group", l)
		}
		hosts = append(hosts, l)
	}
	if groupConfig != "" {
		if hosts, ok := groups[groupConfig]; ok {
			return hosts, nil
		}
		return nil, fmt.Errorf("Group '%s' not found", groupConfig)
	}

	if len(groups) > 0 {
		return nil, fmt.Errorf("Config use group mode but no group specified")
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
