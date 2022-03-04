package scan

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

type Command struct {
	Name string
	Args []string
	cmd *exec.Cmd
}

func (command *Command) run() {
	command.cmd = exec.Command(command.Name, command.Args...)
}

func (command *Command) Start() {
	command.cmd.Start()
}

func (command *Command) Wait() {
	command.cmd.Wait()
}

func (command *Command) NewScanner() *bufio.Scanner {
	command.run()

	stderr, err := command.cmd.StderrPipe()
	if err != nil {
		log.Fatalf("could not get stderr pipe: %v", err)
	}

	stdout, err := command.cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("could not get stdout pipe: %v", err)
	}
	merged := io.MultiReader(stderr, stdout)
	return bufio.NewScanner(merged)
}

func (command *Command) ExitCode() int {
	return command.cmd.ProcessState.ExitCode()
}
