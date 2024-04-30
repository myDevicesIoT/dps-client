package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	MaxExecutionDuration = time.Second * 15
)

var (
	mux sync.RWMutex
)

type GatewayCommandExecRequest struct {
	// Command to execute.
	// This command must be pre-configured in the LoRa Gateway Bridge configuration.
	Command string
	// Execution request ID (UUID).
	// The same token will be returned when the execution of the command has
	// completed.
	ExecId []byte
	// Standard input.
	Stdin []byte
	// Environment variables.
	Environment map[string]string
}

type GatewayCommandExecResponse struct {
	// Gateway ID.
	GatewayId []byte
	// Execution request ID (UUID).
	ExecId []byte
	// Standard output.
	Stdout []byte
	// Standard error.
	Stderr []byte
	// Error message in case the command execution failed.
	Error string
}

func ExecuteCommand(cmd GatewayCommandExecRequest) (GatewayCommandExecResponse, error) {
	stdout, stderr, err := execute(cmd.Command, cmd.Stdin, cmd.Environment)
	resp := GatewayCommandExecResponse{
		ExecId: cmd.ExecId,
		Stdout: stdout,
		Stderr: stderr,
	}
	if err != nil {
		return resp, err
	}

	return resp, nil

}

func execute(command string, stdin []byte, environment map[string]string) ([]byte, []byte, error) {
	mux.RLock()
	defer mux.RUnlock()

	cmdArgs, err := ParseCommandLine(command)
	if err != nil {
		return nil, nil, errors.Wrap(err, "parse command error")
	}
	if len(cmdArgs) == 0 {
		return nil, nil, errors.New("no command is given")
	}

	log.WithFields(log.Fields{
		"command": command,
		"exec":    cmdArgs[0],
		"args":    cmdArgs[1:],
	}).Info("commands: executing command")

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(MaxExecutionDuration))
	defer cancel()

	cmdCtx := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)

	// The default is that when cmdCtx.Env is nil, os.Environ() are being used
	// automatically. As we want to add additional env. variables, we want to
	// extend this list, thus first need to set them to os.Environ()
	cmdCtx.Env = os.Environ()
	for k, v := range environment {
		cmdCtx.Env = append(cmdCtx.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdinPipe, err := cmdCtx.StdinPipe()
	if err != nil {
		return nil, nil, errors.Wrap(err, "get stdin pipe error")
	}

	stdoutPipe, err := cmdCtx.StdoutPipe()
	if err != nil {
		return nil, nil, errors.Wrap(err, "get stdout pipe error")
	}

	stderrPipe, err := cmdCtx.StderrPipe()
	if err != nil {
		return nil, nil, errors.Wrap(err, "get stderr pipe error")
	}

	go func() {
		defer stdinPipe.Close()
		if _, err := stdinPipe.Write(stdin); err != nil {
			log.WithError(err).Error("commands: write to stdin error")
		}
	}()

	if err := cmdCtx.Start(); err != nil {
		return nil, nil, errors.Wrap(err, "starting command error")
	}

	stdoutB, _ := ioutil.ReadAll(stdoutPipe)
	stderrB, _ := ioutil.ReadAll(stderrPipe)

	if err := cmdCtx.Wait(); err != nil {
		return nil, nil, errors.Wrap(err, "waiting for command to finish error")
	}

	return stdoutB, stderrB, nil
}

// ParseCommandLine parses the given command to commands and arguments.
// source: https://stackoverflow.com/questions/34118732/parse-a-command-line-string-into-flags-and-arguments-in-golang
func ParseCommandLine(command string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(command); i++ {
		c := command[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New(fmt.Sprintf("Unclosed quote in command line: %s", command))
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}
