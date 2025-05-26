// Package shell provides utilities for executing shell commands
// and handling their output in a safe and controlled manner.
package shell

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// ExecCommand executes a shell command with the given arguments and returns its output.
// It captures both stdout and stderr, returning stdout as a string on success.
// If the command fails, it returns an error containing both the error message and stderr output.
//
// Parameters:
//   - command: The command to execute
//   - arg: Variable number of command arguments
//
// Returns:
//   - string: The command's stdout output
//   - error: Error if the command fails, including both error message and stderr output
func ExecCommand(command string, arg ...string) (string, error) {
	cmd := exec.Command(command, arg...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		var errStr = fmt.Sprintf(`Command exec error: [%s %s]. %s %s`,
			command, strings.Join(arg, " "), err, stderr.String())

		return "", errors.New(errStr)
	}

	return strings.Trim(out.String(), "\n"), nil
}
