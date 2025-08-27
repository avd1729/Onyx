package runtimes

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os/exec"
	"sandbox/pkg/model"
	"sandbox/pkg/utils"
	"strings"
	"time"
)

type PythonExecutor struct{}

func (p PythonExecutor) Execute(
	ctx context.Context,
	code string,
	timeout time.Duration,
	dependencies ...string,
) model.ExecResult {

	logPrefix := "[sandbox-exec]"
	log.Println(logPrefix, "Preparing to execute Python code in Docker...")

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if !utils.IsDockerAvailable() {
		log.Println(logPrefix, "Docker is not available or not running")
		return model.ExecResult{Err: errors.New("docker is not available or not running")}
	}

	// If dependencies exist, prepare install command
	installCmd := ""
	if len(dependencies) > 0 {
		installCmd = "pip install " + strings.Join(dependencies, " ") + " && "
	}

	dockerCmd := []string{
		"run", "--rm", "-i",
		"python:3.11",
		"sh", "-c", installCmd + "python -",
	}

	log.Println(logPrefix, "Running:", "docker", dockerCmd)

	cmd := exec.CommandContext(ctx, "docker", dockerCmd...)
	cmd.Stdin = bytes.NewReader([]byte(code))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		log.Println(logPrefix, "Execution timed out")
		return model.ExecResult{Err: errors.New("execution timed out")}
	}

	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n" + stderr.String()
		} else {
			output = stderr.String()
		}
	}

	if err != nil {
		log.Printf("%s Docker execution failed: %v\n", logPrefix, err)
		return model.ExecResult{Output: output, Err: errors.New("docker execution failed")}
	}

	log.Println(logPrefix, "Execution succeeded.")
	return model.ExecResult{Output: output}
}
