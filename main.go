package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func main() {
	flag.Parse()
	os.Exit(Main())
}

var (
	colorList = []func(string, ...interface{}) string{
		color.BlueString,
		color.CyanString,
		color.GreenString,
		color.MagentaString,
		color.YellowString,
	}
	colorCounter int
)

var (
	host = flag.String("h", "", "comma separated host")
	edit = flag.Bool("e", false, "edit config")
)

func Main() int {

	if len(os.Args) < 2 {
		fmt.Println("usage: rmt [-h <hosts>|-e] <command>")
		return 0
	}

	var err error
	var hosts []string
	if *host != "" {
		hosts = strings.Split(*host, ",")
	}
	args := flag.Args()

	if *edit {
		err := editConfig()
		if err != nil {
			fmt.Println(errColor(err.Error()))
		}
	}

	if len(hosts) == 0 {
		hosts, err = readHostConfig()
		if err != nil {
			fmt.Println(errColor(err.Error()))
		}
	}

	if len(args) == 0 {
		return 0
	}

	var grCount int
	errChan := make(chan error)
	for _, host := range hosts {
		go run(host, args, errChan)
		grCount++
	}
	for grCount != 0 {
		err := <-errChan
		if err != nil {
			println(err.Error())
		}
		grCount--
	}

	return 0
}

func run(server string, command []string, err chan error) {
	cmds := []string{"ssh", server, strings.Join(command, " ")}
	fmt.Println("Executing : ", cmds)

	cmd := exec.Command(cmds[0], cmds[1:]...)

	cmd.Stdout = newWriter(randColor(fmt.Sprintf("[%s] ", server)))
	cmd.Stderr = newWriter(errColor(fmt.Sprintf("[%s] ERR : ", server)))

	if errno := cmd.Run(); errno != nil {
		err <- fmt.Errorf("[%s] %s", server, errno.Error())
		return
	}

	err <- fmt.Errorf("[%s] %s", server, "session closed")
}

func randColor(s string) string {
	colorCounter++
	if colorCounter == len(colorList) {
		colorCounter = 0
	}
	return colorList[colorCounter](s)
}
func errColor(s string) string {
	return color.RedString(s)
}

type writer struct {
	prefix string
	pipe   chan string
}

func newWriter(prefix string) *writer {
	w := &writer{
		prefix: prefix,
		pipe:   make(chan string),
	}
	go w.run()
	return w
}

func (c *writer) run() {
	for {
		fmt.Print(<-c.pipe)
	}
}

func (c *writer) Write(b []byte) (int, error) {
	c.pipe <- c.prefix + string(b)
	return len(b), nil
}
