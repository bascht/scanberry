package scan

import (
	"strings"
	"testing"
)

func TestCommandRun(t *testing.T) {
	command := Command{Name: "echo", Args: []string{"Here", "be", "a", "test"}}

	scanner := command.NewScanner()
	command.Start()
	for scanner.Scan() {
		text := scanner.Text()
		if text != "Here be a test" {
			t.Error("Didn't do shit with" + text)
		}
	}
}
func TestCommandOutput(t *testing.T) {
	command := Command{Name: "go", Args: []string{"version"}}

	scanner := command.NewScanner()
	command.Start()
	for scanner.Scan() {
		text := scanner.Text()
		if !strings.HasPrefix(text, "go version") {
			t.Error("Didn't do shit with" + text)
		} else {
			t.Logf("Found correct output of %v", text)
		}
	}
	t.Logf("COMMAND %v", command.cmd.ProcessState.ExitCode())

	// TODO: Find out why this is always -1
	// if command.ExitCode() != 0 {
	// 	t.Error("Didn't run successfully", command.ExitCode())
	// }
}

func TestCommandFailure(t *testing.T) {
	command := Command{Name: "go", Args: []string{"gahdned"}}

	scanner := command.NewScanner()
	command.Start()
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	if !strings.HasPrefix(text[0], "go gahdned") {
		t.Error("Didn't do shit with" + text[0])
	}
	t.Logf("## Command Output %v", command.cmd.ProcessState.ExitCode())
}
