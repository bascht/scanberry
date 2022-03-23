package scan

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

type Command struct {
	Name string
	Args []string
	cmd *exec.Cmd
}

func (command *Command) GetScanner() *bufio.Scanner {
	command.cmd = exec.Command(command.Name, command.Args...)

	stdout, err := command.cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("could not get stdout pipe: %v", err)
	}
	command.cmd.Stderr = command.cmd.Stdout

	return bufio.NewScanner(stdout)
}

func (command *Command) Start() {
	fmt.Sprintf("Running Command %v \n", command)
	command.cmd.Start()
}

func (command *Command) Wait() {
	command.cmd.Wait()
}

func (command *Command) RanSuccessful() bool {
	return command.cmd.ProcessState.ExitCode() == 0
}
